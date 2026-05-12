// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds query DTOs for the time-windowed feed and chart
// endpoints. Different endpoints disagree on whether timestamps are
// in milliseconds or seconds — the field names below mirror the
// Derive docs so the unit is obvious at the call site.
package types

// FundingRateHistoryQuery selects funding-rate prints for one
// perpetual instrument. Timestamps are in milliseconds since the
// Unix epoch (zero values defer to the server defaults).
type FundingRateHistoryQuery struct {
	HistoryWindow
	// InstrumentName is the perpetual to fetch rates for. Required.
	InstrumentName string
	// Period optionally narrows the response to one period bucket.
	Period string
}

// Validate performs schema-level checks on the receiver.
func (q FundingRateHistoryQuery) Validate() error {
	if q.InstrumentName == "" {
		return invalidParam("instrument_name", "required")
	}
	return nil
}

// SpotFeedHistoryQuery selects oracle spot-feed prints for one
// currency. PeriodSec accepts one of the engine's bucket sizes (60,
// 300, 900, 1800, 3600, 14400, 28800, 86400, 604800). Timestamps
// are in milliseconds since the Unix epoch.
type SpotFeedHistoryQuery struct {
	HistoryWindow
	// Currency is the asset symbol (e.g. "BTC"). Required.
	Currency string
	// PeriodSec is the bucket size in seconds. Required.
	PeriodSec int64
}

// Validate performs schema-level checks on the receiver.
func (q SpotFeedHistoryQuery) Validate() error {
	if q.Currency == "" {
		return invalidParam("currency", "required")
	}
	if q.PeriodSec <= 0 {
		return invalidParam("period", "required")
	}
	return nil
}

// SpotFeedHistoryCandlesQuery selects OHLC candles for one
// currency's spot feed. Timestamps are in milliseconds since the
// Unix epoch.
type SpotFeedHistoryCandlesQuery struct {
	HistoryWindow
	// Currency is the asset symbol. Required.
	Currency string
	// PeriodSec is the candle period in seconds. Required.
	PeriodSec int64
}

// Validate performs schema-level checks on the receiver.
func (q SpotFeedHistoryCandlesQuery) Validate() error {
	if q.Currency == "" {
		return invalidParam("currency", "required")
	}
	if q.PeriodSec <= 0 {
		return invalidParam("period", "required")
	}
	return nil
}

// IndexChartQuery selects OHLC candles for one currency's index
// feed. Timestamps are in milliseconds since the Unix epoch.
type IndexChartQuery struct {
	HistoryWindow
	// Currency is the asset symbol. Required.
	Currency string
	// PeriodSec is the candle period in seconds. Required.
	PeriodSec int64
}

// Validate performs schema-level checks on the receiver.
func (q IndexChartQuery) Validate() error {
	if q.Currency == "" {
		return invalidParam("currency", "required")
	}
	if q.PeriodSec <= 0 {
		return invalidParam("period", "required")
	}
	return nil
}

// TradingViewChartQuery selects TradingView UDF-format OHLC bars
// for one instrument. Timestamps are in milliseconds since the Unix
// epoch.
type TradingViewChartQuery struct {
	HistoryWindow
	// InstrumentName is the market to fetch bars for. Required.
	InstrumentName string
	// PeriodSec is the candle period in seconds. Required.
	PeriodSec int64
}

// Validate performs schema-level checks on the receiver.
func (q TradingViewChartQuery) Validate() error {
	if q.InstrumentName == "" {
		return invalidParam("instrument_name", "required")
	}
	if q.PeriodSec <= 0 {
		return invalidParam("period", "required")
	}
	return nil
}

// InterestRateHistoryQuery selects historical USDC borrow / supply
// APY prints. Timestamps are in seconds since the Unix epoch — note
// the unit difference from the other history queries.
type InterestRateHistoryQuery struct {
	// FromSec is the start of the window in seconds since the Unix
	// epoch. Required.
	FromSec int64
	// ToSec is the end of the window in seconds since the Unix
	// epoch. Required.
	ToSec int64
}

// Validate performs schema-level checks on the receiver.
func (q InterestRateHistoryQuery) Validate() error {
	if q.FromSec <= 0 {
		return invalidParam("from_timestamp_sec", "required")
	}
	if q.ToSec <= 0 {
		return invalidParam("to_timestamp_sec", "required")
	}
	return nil
}

// STDRVSnapshotsQuery selects one wallet's staked-DRV balance
// snapshots over a time window. Timestamps are in seconds since the
// Unix epoch.
type STDRVSnapshotsQuery struct {
	// Wallet is the account address. Required.
	Wallet string
	// FromSec optionally narrows the window's start in seconds.
	FromSec int64
	// ToSec optionally narrows the window's end in seconds.
	ToSec int64
}

// Validate performs schema-level checks on the receiver.
func (q STDRVSnapshotsQuery) Validate() error {
	if q.Wallet == "" {
		return invalidParam("wallet", "required")
	}
	return nil
}
