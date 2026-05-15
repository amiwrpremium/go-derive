// Package ws is the WebSocket-backed client for Derive's JSON-RPC API.
//
// This file declares the documented enum values for the channels that
// take a typed parameter on the wire. Use them with [Client.SubscribeOrderBook]
// (group/depth), [Client.SubscribeTicker] / [Client.SubscribeTickerSlim]
// (interval) — they surface in IDE autocompletion and prevent typos
// the docstring otherwise wouldn't catch.
package ws

// Orderbook-channel depth values (number of price levels per side).
// Mirror the enum at
// https://docs.derive.xyz/reference/orderbook-instrument_name-group-depth.
const (
	// Depth1 is top-of-book only.
	Depth1 int = 1
	// Depth10 is the default for [Client.SubscribeOrderBook] when the
	// caller passes zero.
	Depth10 int = 10
	// Depth20 is the next tier up.
	Depth20 int = 20
	// Depth100 is the deepest documented level.
	Depth100 int = 100
	// DepthDefault is the value [Client.SubscribeOrderBook] uses when
	// the caller passes zero.
	DepthDefault = Depth10
)

// Orderbook-channel group values (price-grouping / rounding bucket
// size). Mirror the enum at
// https://docs.derive.xyz/reference/orderbook-instrument_name-group-depth.
const (
	// Group1 disables price grouping (every distinct price level is
	// reported on its own line).
	Group1 string = "1"
	// Group10 buckets adjacent levels in groups of 10 ticks.
	Group10 string = "10"
	// Group100 buckets in groups of 100 ticks.
	Group100 string = "100"
	// GroupDefault is the value [Client.SubscribeOrderBook] uses when
	// the caller passes the empty string.
	GroupDefault = Group1
)

// Ticker-channel interval values (snapshot emission cadence, in
// milliseconds). Mirror the enum at
// https://docs.derive.xyz/reference/ticker-instrument_name-interval.
const (
	// Interval100 emits at ~100 ms — every change, near-realtime.
	Interval100 string = "100"
	// Interval1000 emits every second — the default for
	// [Client.SubscribeTicker] and [Client.SubscribeTickerSlim].
	Interval1000 string = "1000"
	// IntervalDefault is the value used when the caller passes the
	// empty string.
	IntervalDefault = Interval1000
)
