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

import "encoding/json"

// Ticker is the public market summary for one instrument: top-of-book, marks,
// and depth at 5%.
type Ticker struct {
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// InstrumentType is "perp", "option", or "erc20".
	InstrumentType string `json:"instrument_type,omitempty"`
	// IsActive reports whether the instrument is currently open for trading.
	IsActive bool `json:"is_active,omitempty"`

	// BestBidPrice is the highest resting bid.
	BestBidPrice Decimal `json:"best_bid_price"`
	// BestBidAmount is the size resting at [BestBidPrice].
	BestBidAmount Decimal `json:"best_bid_amount"`
	// BestAskPrice is the lowest resting ask.
	BestAskPrice Decimal `json:"best_ask_price"`
	// BestAskAmount is the size resting at [BestAskPrice].
	BestAskAmount Decimal `json:"best_ask_amount"`

	// FivePercentBidDepth is the cumulative bid size within 5 % of mark.
	FivePercentBidDepth Decimal `json:"five_percent_bid_depth,omitempty"`
	// FivePercentAskDepth is the cumulative ask size within 5 % of mark.
	FivePercentAskDepth Decimal `json:"five_percent_ask_depth,omitempty"`

	// MarkPrice is the engine's mark price for the instrument.
	MarkPrice Decimal `json:"mark_price"`
	// IndexPrice is the underlying index price.
	IndexPrice Decimal `json:"index_price"`
	// MinPrice is the engine-enforced lower price band.
	MinPrice Decimal `json:"min_price,omitempty"`
	// MaxPrice is the engine-enforced upper price band.
	MaxPrice Decimal `json:"max_price,omitempty"`

	// OpenInterest is preserved as raw JSON because Derive returns it as
	// a per-margin-type breakdown
	// (`{"PM": [...], "PM2": [...], "SM": [...]}` of `{current_open_interest,
	// interest_cap, manager_currency}` items). Decode further if needed.
	OpenInterest json.RawMessage `json:"open_interest,omitempty"`

	// Timestamp is when this ticker snapshot was produced.
	Timestamp MillisTime `json:"timestamp"`
}
