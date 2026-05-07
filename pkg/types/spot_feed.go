// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// All numeric fields use [Decimal], a thin wrapper around shopspring/decimal,
// so price/size/fee values never lose precision through float64 round-trips.
// On the wire, [Decimal] reads and writes JSON strings (Derive's preferred
// representation); a fallback path also accepts JSON numbers for resilience.
//
// Identifier types ([Address], [TxHash], [MillisTime]) carry the same
// round-trip guarantees: each one preserves the canonical wire format
// regardless of how Go marshals the surrounding struct.
//
// # Why named types
//
// Plain string and int64 fields would parse just fine, but named types let
// the SDK enforce invariants at construction time (NewAddress checksum
// check, NewDecimal precision check) and let callers tell at a glance which
// values are amounts vs prices vs subaccount ids.
package types

// SpotFeed is one update from the `spot_feed.{currency}` WebSocket channel.
//
// Derive's oracle feed delivers a per-currency snapshot of mark and 24h-prior
// prices. Use it for liquidation monitoring, basis calculations, or any
// risk surface that needs an oracle anchor independent of the order book.
type SpotFeed struct {
	// Timestamp is the message-emission time in milliseconds.
	Timestamp MillisTime `json:"timestamp"`
	// Feeds is keyed by currency symbol (e.g. "BTC", "ETH") and carries
	// one [SpotFeedEntry] per currency reported in the message. With the
	// per-currency subscription pattern there is usually exactly one entry.
	Feeds map[string]SpotFeedEntry `json:"feeds"`
}

// SpotFeedEntry is the per-currency oracle reading inside a [SpotFeed].
type SpotFeedEntry struct {
	// Price is the current oracle price (decimal string in quote units).
	Price Decimal `json:"price"`
	// Confidence is the oracle confidence score (decimal string in [0, 1]).
	Confidence Decimal `json:"confidence"`
	// PricePrevDaily is the price 24 hours prior, used to derive the
	// 24h delta on the UI without an extra round trip.
	PricePrevDaily Decimal `json:"price_prev_daily"`
	// ConfidencePrevDaily is the oracle confidence 24 hours prior.
	ConfidencePrevDaily Decimal `json:"confidence_prev_daily"`
	// TimestampPrevDaily is the millisecond timestamp of the 24h-prior
	// reading.
	TimestampPrevDaily MillisTime `json:"timestamp_prev_daily"`
}
