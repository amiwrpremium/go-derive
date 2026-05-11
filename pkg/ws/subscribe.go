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
	out := &Subscription[T]{
		raw:     sub,
		typed:   make(chan T, cfg.bufferSize),
		channel: channelName,
		cfg:     cfg,
	}
	go out.pump()
	return out, nil
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
// subscription. The zero value is not usable; obtain one from [Subscribe].
//
// Always call [Subscription.Close] to release the channel slot and tell
// the server to stop sending updates. The Close call is idempotent.
type Subscription[T any] struct {
	raw     transport.Subscription
	typed   chan T
	channel string
	cfg     subscribeConfig
}

// Channel returns the dotted server-side channel name (e.g.
// "orderbook.BTC-PERP.1.10"). Useful for diagnostics and logs.
func (s *Subscription[T]) Channel() string { return s.channel }

// Updates returns the receive channel of typed events. The channel is
// closed when the subscription terminates; receivers should select against
// ctx.Done() to know when to bail.
func (s *Subscription[T]) Updates() <-chan T { return s.typed }

// Err returns the terminal error once [Subscription.Updates] has closed,
// or nil for a clean shutdown.
func (s *Subscription[T]) Err() error { return s.raw.Err() }

// Close ends the subscription, sends an unsubscribe RPC best-effort, and
// drains the typed channel. Idempotent.
func (s *Subscription[T]) Close() error { return s.raw.Close() }

// pump bridges the untyped transport channel to the typed user channel,
// applying the configured drop policy when the buffer is full.
func (s *Subscription[T]) pump() {
	defer close(s.typed)
	for v := range s.raw.Updates() {
		typed, ok := v.(T)
		if !ok {
			s.notify(ErrTypeMismatch)
			continue
		}
		s.deliver(typed)
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
