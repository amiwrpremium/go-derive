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
// The returned subscription buffers up to 256 events in memory; if the
// caller is slow, newer events are dropped (best-effort fan-out, not a
// reliable queue). Use [SubscribeFunc] when you want to be sure every event
// is processed.
func Subscribe[T any](ctx context.Context, c *Client, channelName string, decoder func(json.RawMessage) (T, error)) (*Subscription[T], error) {
	dec := func(raw json.RawMessage) (any, error) {
		return decoder(raw)
	}
	sub, err := c.transport().Subscribe(ctx, channelName, dec)
	if err != nil {
		return nil, err
	}
	out := &Subscription[T]{
		raw:     sub,
		typed:   make(chan T, 256),
		channel: channelName,
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
func SubscribeFunc[T any](ctx context.Context, c *Client, channelName string, decoder func(json.RawMessage) (T, error), fn func(T)) error {
	sub, err := Subscribe[T](ctx, c, channelName, decoder)
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
