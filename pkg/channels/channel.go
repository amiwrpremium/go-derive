// Package channels declares the typed WebSocket subscription channels
// Derive supports.
//
// Each channel descriptor (see pkg/channels/public and pkg/channels/private)
// implements [Channel]: it knows the dotted server-side name to subscribe
// to and how to decode an inbound notification payload into a typed Go
// value. The descriptors are passed into
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe] and
// [github.com/amiwrpremium/go-derive/pkg/ws.SubscribeFunc].
//
// # Public vs private
//
// Public channels (pkg/channels/public) need no authentication. Private
// channels (pkg/channels/private) need a logged-in WebSocket — the
// SubaccountID field on each private descriptor scopes the stream to
// one subaccount.
package channels

import "encoding/json"

// Channel is the contract a subscription descriptor must satisfy.
//
// Implementations live in pkg/channels/public and pkg/channels/private.
// Third-party implementations are welcome — anything that produces a
// dotted name and decodes JSON to a typed value will work.
type Channel interface {
	// Name returns the dot-separated server-side channel name, e.g.
	// "trades.BTC-PERP" or "subaccount.123.orders".
	Name() string

	// Decode turns a raw notification payload into a typed value. The
	// concrete return type is descriptor-specific; pass the matching T to
	// ws.Subscribe to consume it without an explicit type assertion.
	Decode(raw json.RawMessage) (any, error)
}
