// Package types — see address.go for the overview.
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
