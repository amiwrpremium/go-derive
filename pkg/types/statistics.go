// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shape of `public/statistics`.
package types

// Statistics is the response of `public/statistics`. The endpoint
// returns rolling 24-hour and all-time statistics for one
// instrument: volume, premium volume, fees, trades count, plus
// total open interest.
//
// Mirrors `PublicStatisticsResultSchema` in Derive's v2.2 OpenAPI
// spec.
type Statistics struct {
	// DailyFees is the 24h fees taken on the instrument.
	DailyFees Decimal `json:"daily_fees"`
	// DailyNotionalVolume is the 24h notional volume.
	DailyNotionalVolume Decimal `json:"daily_notional_volume"`
	// DailyPremiumVolume is the 24h premium volume.
	DailyPremiumVolume Decimal `json:"daily_premium_volume"`
	// DailyTrades is the 24h trade count.
	DailyTrades int64 `json:"daily_trades"`
	// OpenInterest is the current open interest on the instrument.
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
