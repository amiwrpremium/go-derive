// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the per-(currency, instrument_type) aggregate
// statistics returned by `public/all_statistics`.
package types

// AggregateStatistics is one entry in `public/all_statistics`.
// Each entry is one (currency, instrument_type) tuple's rolling 24h
// and all-time statistics — the per-instrument [Statistics] shape
// rolled up across every instrument matching the tuple.
//
// Mirrors the response shape per
// docs.derive.xyz/reference/public-all_statistics.
type AggregateStatistics struct {
	// Currency is the underlying currency (e.g. "ETH").
	Currency string `json:"currency"`
	// InstrumentType is the instrument class (e.g. "perp",
	// "option", "erc20").
	InstrumentType string `json:"instrument_type"`
	// DailyFees is the rolling 24h fees taken on the tuple.
	DailyFees Decimal `json:"daily_fees"`
	// DailyNotionalVolume is the rolling 24h notional volume.
	DailyNotionalVolume Decimal `json:"daily_notional_volume"`
	// DailyPremiumVolume is the rolling 24h premium volume.
	DailyPremiumVolume Decimal `json:"daily_premium_volume"`
	// DailyTrades is the rolling 24h trade count.
	DailyTrades int64 `json:"daily_trades"`
	// OpenInterest is the current total open interest on the tuple.
	OpenInterest Decimal `json:"open_interest"`
	// TotalFees is the all-time fees taken.
	TotalFees Decimal `json:"total_fees"`
	// TotalNotionalVolume is the all-time notional volume.
	TotalNotionalVolume Decimal `json:"total_notional_volume"`
	// TotalPremiumVolume is the all-time premium volume.
	TotalPremiumVolume Decimal `json:"total_premium_volume"`
	// TotalTrades is the all-time trade count.
	TotalTrades int64 `json:"total_trades"`
}
