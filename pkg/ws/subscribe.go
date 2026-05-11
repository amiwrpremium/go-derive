// Package ws is the WebSocket-backed client for Derive's JSON-RPC API.
//
// # What it covers
//
// Derive's WebSocket transport carries two distinct workloads:
//
//   - request/response RPCs (lower latency than HTTP because of connection
//     reuse and no per-call TLS handshake)
//   - pub/sub channel notifications (the only way to stream live data)
//
// [Client] handles both. It runs three goroutines under one parent context:
// a read pump that demultiplexes responses from notifications, a write pump
// that serialises outgoing frames, and a ping pump that keeps the connection
// alive. When [WithReconnect] is enabled, a reconnect goroutine re-dials
// with exponential backoff and re-issues subscribe + login on success.
//
// # Lifecycle
//
//	c, _ := ws.New(ws.WithMainnet(), ws.WithSigner(s), ws.WithSubaccount(123))
//	defer c.Close()
//	if err := c.Connect(ctx); err != nil { ... }
//	if err := c.Login(ctx); err != nil { ... }
//
//	sub, err := c.SubscribeOrderBook(ctx, "BTC-PERP", "", 0)
//	defer sub.Close()
//	for ob := range sub.Updates() { ... }
//
// # Concurrency
//
// [Client] is safe for concurrent use after Connect. Many goroutines may
// call methods or hold subscriptions on the same client simultaneously.
package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"sync"

	"github.com/amiwrpremium/go-derive/internal/transport"
)

// Subscribe registers a typed subscription on a [Client] and returns a
// [Subscription] whose Updates channel yields values of type T.
//
// channelName is the dotted wire string from the docs (e.g.
// "orderbook.BTC-PERP.1.10" or "7.orders"). decoder turns one
// notification's raw payload into a typed T; for the common case where
// payloads are plain JSON, pass [decodeJSON][T] (the convenience the
// typed Subscribe* methods on [Client] use internally).
//
// Most callers should prefer the typed methods on [Client] —
// [Client.SubscribeOrderBook], [Client.SubscribeOrders], etc. They
// bake in the right name and decoder so you only pass the channel
// parameters. Use this generic form when you need a channel that
// isn't documented yet, or to attach a custom decoder.
//
// Generics let callers avoid type assertions at the use site:
//
//	sub, _ := ws.Subscribe(ctx, c, "orderbook.BTC-PERP.1.10",
//	    func(raw json.RawMessage) (types.OrderBook, error) {
//	        var ob types.OrderBook
//	        return ob, json.Unmarshal(raw, &ob)
//	    })
//	defer sub.Close()
//	for ob := range sub.Updates() {
//	    fmt.Println(ob.Bids[0])
//	}
//
// # Lifetime
//
// ctx controls the subscription's lifetime: cancelling it tears
// the subscription down (sends unsubscribe, closes Updates, exits
// the pump goroutines). Pass [context.Background] to keep the
// subscription alive until explicit [Subscription.Close].
//
// The returned subscription buffers up to 256 events in memory by
// default; if the caller is slow, the [DropNewest] policy silently
// drops incoming events (the default). Tune with [WithBufferSize],
// [WithDropPolicy], and [WithErrorHandler]. Use [SubscribeFunc] when
// you want to drive a callback that back-pressures naturally.
func Subscribe[T any](ctx context.Context, c *Client, channelName string, decoder func(json.RawMessage) (T, error), opts ...SubscribeOption) (*Subscription[T], error) {
	cfg := applySubscribeOpts(opts)
	dec := func(raw json.RawMessage) (any, error) {
		return decoder(raw)
	}
	sub, err := c.transport().Subscribe(ctx, channelName, dec)
	if err != nil {
		return nil, err
	}
	out := newSubscription[T](sub, make(chan T, cfg.bufferSize), channelName, cfg, true /* ownsTyped */)
	out.start(ctx)
	return out, nil
}

// SubscribeInto registers a typed subscription that delivers events
// into the caller-supplied out channel. Use it when fanning multiple
// subscriptions into one shared consumer loop, or when the caller
// already owns a sized channel they want filled.
//
// The caller owns out: the SDK never closes it. Close the
// subscription with [Subscription.Close] when done; only after that
// call returns is it safe for the caller to close out.
//
// Apart from buffer ownership the semantics match [Subscribe]:
// drop policy and error handler still apply via [SubscribeOption].
// [WithBufferSize] has no effect — the caller's chan determines
// capacity. ctx controls the subscription's lifetime (see the
// "Lifetime" section on [Subscribe]).
func SubscribeInto[T any](ctx context.Context, c *Client, channelName string, decoder func(json.RawMessage) (T, error), out chan T, opts ...SubscribeOption) (*Subscription[T], error) {
	cfg := applySubscribeOpts(opts)
	dec := func(raw json.RawMessage) (any, error) {
		return decoder(raw)
	}
	sub, err := c.transport().Subscribe(ctx, channelName, dec)
	if err != nil {
		return nil, err
	}
	s := newSubscription[T](sub, out, channelName, cfg, false /* ownsTyped */)
	s.start(ctx)
	return s, nil
}

// SubscribeFunc is a convenience over [Subscribe] that drives a per-event
// callback synchronously. It returns when ctx is cancelled (returning
// ctx.Err()) or the subscription closes (returning the underlying
// terminal error, which may be nil for a clean close).
//
// Use SubscribeFunc when callback-driven code reads more naturally than a
// channel-receive loop, or when you want to guarantee every event is
// processed (the callback runs synchronously, so back-pressure on the
// caller is back-pressure on the subscription).
func SubscribeFunc[T any](ctx context.Context, c *Client, channelName string, decoder func(json.RawMessage) (T, error), fn func(T), opts ...SubscribeOption) error {
	sub, err := Subscribe[T](ctx, c, channelName, decoder, opts...)
	if err != nil {
		return err
	}
	defer func() { _ = sub.Close() }()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-sub.Updates():
			if !ok {
				return sub.Err()
			}
			fn(ev)
		}
	}
}

// Subscription is a typed wrapper around the underlying transport-level
// subscription. The zero value is not usable; obtain one from [Subscribe]
// or [SubscribeInto].
//
// The subscription terminates when any of these happens:
//   - the context passed to Subscribe is cancelled,
//   - [Subscription.Close] is called explicitly,
//   - the WebSocket connection drops (with [WithReconnect] disabled).
//
// [Subscription.Close] is idempotent.
type Subscription[T any] struct {
	raw       transport.Subscription
	typed     chan T
	channel   string
	cfg       subscribeConfig
	ownsTyped bool // true for Subscribe; false for SubscribeInto (caller owns the chan)

	done     chan struct{} // closed once when the subscription terminates for any reason
	doneOnce sync.Once
}

// newSubscription constructs a Subscription[T] with the done-channel
// machinery wired up. Callers should immediately call start to spawn
// the pump / drainErrors / watchCtx goroutines.
func newSubscription[T any](raw transport.Subscription, typed chan T, channel string, cfg subscribeConfig, ownsTyped bool) *Subscription[T] {
	return &Subscription[T]{
		raw:       raw,
		typed:     typed,
		channel:   channel,
		cfg:       cfg,
		ownsTyped: ownsTyped,
		done:      make(chan struct{}),
	}
}

// start spawns the three lifecycle goroutines: pump (typed delivery),
// drainErrors (decode-error forwarding), and watchCtx (lifetime hook).
func (s *Subscription[T]) start(ctx context.Context) {
	go s.pump()
	go s.drainErrors()
	go s.watchCtx(ctx)
}

// Channel returns the dotted server-side channel name (e.g.
// "orderbook.BTC-PERP.1.10"). Useful for diagnostics and logs.
func (s *Subscription[T]) Channel() string { return s.channel }

// Updates returns the receive channel of typed events. The channel
// is closed when the subscription terminates — either because the
// Subscribe ctx was cancelled, [Subscription.Close] was called, or
// the WebSocket connection dropped. Receivers can still select
// against ctx.Done() if they want to, but it is no longer strictly
// necessary: cancelling the Subscribe ctx closes Updates directly.
//
// # Multi-subscription gotcha
//
// Each Subscription has its own buffer and drop policy, independent of
// every other Subscription on the same Client. When you select across
// multiple Subscriptions in one goroutine and one handler arm is slow,
// the other subs' buffers keep filling while you're blocked — and
// once full, those other subs start dropping per their configured
// DropPolicy. The drop happens on the SLOW path's neighbours, not the
// slow path itself.
//
// Three ways out:
//
//   - spawn a goroutine per Subscription so each handler runs
//     independently (the simplest and usually right answer);
//   - register [WithErrorHandler] so [ErrBufferFull] is observable
//     when it happens;
//   - or fan-in via [SubscribeInto] with one shared chan when
//     shared back-pressure across subs is what you actually want.
func (s *Subscription[T]) Updates() <-chan T { return s.typed }

// Err returns the terminal error once [Subscription.Updates] has closed,
// or nil for a clean shutdown (including a ctx-cancel-driven teardown).
func (s *Subscription[T]) Err() error { return s.raw.Err() }

// All returns an [iter.Seq2] over the subscription's events suitable
// for use with Go 1.23+ range-over-func:
//
//	for ev, err := range sub.All() {
//	    if err != nil {
//	        return err
//	    }
//	    handle(ev)
//	}
//
// The first value is the typed event; the second is non-nil only on
// the final terminal yield if the subscription ended in error.
// Clean shutdowns (ctx-cancel, explicit [Subscription.Close]) end
// the iterator without yielding a terminal error.
//
// The iterator stops when the loop body breaks (yield returns false),
// the subscription terminates, or the buffer is drained after
// termination.
//
// Decode errors during the subscription do NOT come through this
// iterator — they are routed through [WithErrorHandler] as
// [ErrDecodeFailed]. Only terminal errors (transport fault) are
// yielded here.
func (s *Subscription[T]) All() iter.Seq2[T, error] {
	if s.ownsTyped {
		// Subscribe path: pump closes typed on exit. A range loop
		// naturally drains any buffered events then signals
		// termination via ok=false, so we don't need to select on
		// done here. Doing so would race — Go picks ready cases
		// at random, so a closed-done arm could win against a
		// non-empty typed buffer and skip events.
		return func(yield func(T, error) bool) {
			for ev := range s.typed {
				if !yield(ev, nil) {
					return
				}
			}
			s.yieldTerminalErr(yield)
		}
	}
	// SubscribeInto path: caller owns typed and the SDK never
	// closes it. Use done as the termination signal, but prefer
	// draining typed first so buffered events aren't lost to the
	// select-arm coin-flip when both arms are ready.
	return func(yield func(T, error) bool) {
		for {
			select {
			case ev := <-s.typed:
				if !yield(ev, nil) {
					return
				}
				continue
			default:
			}
			// typed empty; wait for the next event or termination.
			select {
			case ev := <-s.typed:
				if !yield(ev, nil) {
					return
				}
			case <-s.done:
				s.yieldTerminalErr(yield)
				return
			}
		}
	}
}

// yieldTerminalErr emits one final (zero, err) pair if the
// subscription ended in a non-nil error. Used by [Subscription.All]
// at iterator-exit time.
func (s *Subscription[T]) yieldTerminalErr(yield func(T, error) bool) {
	if err := s.Err(); err != nil {
		var zero T
		yield(zero, err)
	}
}

// Close ends the subscription, sends an unsubscribe RPC best-effort, and
// signals the pump goroutines to exit. Idempotent.
func (s *Subscription[T]) Close() error {
	s.markDone()
	return s.raw.Close()
}

// markDone closes the done channel exactly once. Safe to call from
// any goroutine.
func (s *Subscription[T]) markDone() {
	s.doneOnce.Do(func() { close(s.done) })
}

// watchCtx ties the subscription's lifetime to ctx. When ctx is
// cancelled it calls Close, which tears down the transport
// subscription and signals every other goroutine to exit. When the
// subscription terminates for another reason (explicit Close or
// transport disconnect), done is already closed and the second arm
// fires, letting this goroutine exit without leaking.
func (s *Subscription[T]) watchCtx(ctx context.Context) {
	select {
	case <-ctx.Done():
		_ = s.Close()
	case <-s.done:
	}
}

// pump bridges the untyped transport channel to the typed user channel,
// applying the configured drop policy when the buffer is full. Only
// closes the typed channel when the SDK owns it ([Subscribe]); for
// [SubscribeInto] the caller-supplied channel is left alone.
//
// Exits on any of: transport closed Updates, [Subscription.Close], or
// ctx cancellation (via watchCtx → Close → done). All three converge
// on the done channel.
func (s *Subscription[T]) pump() {
	defer s.markDone()
	if s.ownsTyped {
		defer close(s.typed)
	}
	for {
		select {
		case v, ok := <-s.raw.Updates():
			if !ok {
				return
			}
			typed, ok := v.(T)
			if !ok {
				s.notify(ErrTypeMismatch)
				continue
			}
			s.deliver(typed)
		case <-s.done:
			return
		}
	}
}

// deliver pushes one event onto the typed channel honoring the
// configured drop policy.
func (s *Subscription[T]) deliver(v T) {
	switch s.cfg.dropPolicy {
	case DropOldest:
		select {
		case s.typed <- v:
		default:
			select {
			case <-s.typed: // pop oldest, best-effort
			default:
			}
			select {
			case s.typed <- v:
			default:
				s.notify(ErrBufferFull)
			}
		}
	case Block:
		s.typed <- v
	default: // DropNewest
		select {
		case s.typed <- v:
		default:
			s.notify(ErrBufferFull)
		}
	}
}

// notify invokes the user-supplied error handler if any. Runs on the
// pump goroutine — handlers should be non-blocking.
func (s *Subscription[T]) notify(err error) {
	if s.cfg.errorHandler != nil {
		s.cfg.errorHandler(err)
	}
}

// drainErrors forwards transport-level decoder errors to the
// configured error handler, wrapped with [ErrDecodeFailed]. Exits
// when the transport closes its errors channel OR the subscription
// terminates (signalled via done).
func (s *Subscription[T]) drainErrors() {
	defer s.markDone()
	for {
		select {
		case err, ok := <-s.raw.Errors():
			if !ok {
				return
			}
			if s.cfg.errorHandler == nil {
				continue
			}
			s.cfg.errorHandler(fmt.Errorf("%w: %w", ErrDecodeFailed, err))
		case <-s.done:
			return
		}
	}
}
