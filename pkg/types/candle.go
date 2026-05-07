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

// Candle is one OHLC bar returned by the trade-history endpoints when a
// time-series view is requested.
type Candle struct {
	// StartTimestamp is the bar's open time.
	StartTimestamp MillisTime `json:"timestamp"`
	// Open is the first traded price in the bar.
	Open Decimal `json:"open"`
	// High is the highest traded price in the bar.
	High Decimal `json:"high"`
	// Low is the lowest traded price in the bar.
	Low Decimal `json:"low"`
	// Close is the last traded price in the bar.
	Close Decimal `json:"close"`
	// Volume is the sum of traded sizes in base-currency units.
	Volume Decimal `json:"volume,omitempty"`
}
