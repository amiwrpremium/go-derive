// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shapes for the public chart-data
// endpoints. The candle shape returned by `public/get_index_chart_data`
// and `public/get_spot_feed_history_candles` is the same — see
// [SpotFeedCandle] in feeds.go.
package types

// TradingViewChart is one bar from
// `public/get_tradingview_chart_data`. The endpoint backs the
// TradingView UDF charting protocol, so each bar carries OHLC plus
// per-bar volume in both contracts and USD.
//
// The shape mirrors `TradingviewChartDataResponseSchema` in Derive's
// v2.2 OpenAPI spec.
type TradingViewChart struct {
	// Timestamp is the bucket-start timestamp the bar represents.
	Timestamp MillisTime `json:"timestamp"`
	// TimestampBucket is the regularly-spaced bucket-start time —
	// always equal to Timestamp on this endpoint, but emitted
	// separately by the OAS for parity with the spot/index endpoints.
	TimestampBucket MillisTime `json:"timestamp_bucket"`
	// OpenPrice is the bucket's open price.
	OpenPrice Decimal `json:"open_price"`
	// HighPrice is the bucket's high price.
	HighPrice Decimal `json:"high_price"`
	// LowPrice is the bucket's low price.
	LowPrice Decimal `json:"low_price"`
	// ClosePrice is the bucket's close price.
	ClosePrice Decimal `json:"close_price"`
	// VolumeContracts is the volume traded in the bucket, in
	// contracts.
	VolumeContracts Decimal `json:"volume_contracts"`
	// VolumeUSD is the notional volume traded in the bucket.
	VolumeUSD Decimal `json:"volume_usd"`
}
