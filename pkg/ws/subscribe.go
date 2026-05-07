// Package ws.
package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/amiwrpremium/go-derive/internal/transport"
	"github.com/amiwrpremium/go-derive/pkg/channels"
)

// Subscribe registers a typed subscription on a [Client] and returns a
// [Subscription] whose Updates channel yields values of type T.
//
// T must match the type the channel descriptor's Decode method returns; a
// mismatch is dropped silently rather than crashing the read pump (the
// underlying decoder error is surfaced if a debugger is attached). Pass
// the right T for the descriptor — e.g. types.OrderBook for
// public.OrderBook, []types.Trade for public.Trades.
//
// Generics let callers avoid type assertions at the use site:
//
//	sub, _ := ws.Subscribe[types.OrderBook](ctx, c,
//	    public.OrderBook{Instrument: "BTC-PERP"})
//	defer sub.Close()
//	for ob := range sub.Updates() {
//	    fmt.Println(ob.Bids[0])
//	}
//
// The returned subscription buffers up to 256 events in memory; if the
// caller is slow, newer events are dropped (best-effort fan-out, not a
// reliable queue). Use [SubscribeFunc] when you want to be sure every event
// is processed.
func Subscribe[T any](ctx context.Context, c *Client, ch channels.Channel) (*Subscription[T], error) {
	dec := func(raw json.RawMessage) (any, error) {
		v, err := ch.Decode(raw)
		if err != nil {
			return nil, err
		}
		typed, ok := v.(T)
		if !ok {
			return nil, fmt.Errorf("ws: channel %q: decoded type %T does not match expected %T", ch.Name(), v, *new(T))
		}
		return typed, nil
	}
	sub, err := c.transport().Subscribe(ctx, ch.Name(), dec)
	if err != nil {
		return nil, err
	}
	out := &Subscription[T]{
		raw:     sub,
		typed:   make(chan T, 256),
		channel: ch.Name(),
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
func SubscribeFunc[T any](ctx context.Context, c *Client, ch channels.Channel, fn func(T)) error {
	sub, err := Subscribe[T](ctx, c, ch)
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

// pump bridges the untyped transport channel to the typed user channel.
// Type-mismatched events are dropped (Subscribe returns an error if T
// can't accept the descriptor's output, but we still defend at runtime).
func (s *Subscription[T]) pump() {
	defer close(s.typed)
	for v := range s.raw.Updates() {
		typed, ok := v.(T)
		if !ok {
			continue
		}
		s.typed <- typed
	}
}
