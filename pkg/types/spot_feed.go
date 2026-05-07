// Package types — see address.go for the overview.
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
