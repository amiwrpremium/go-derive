// Package types.
package types

import "encoding/json"

// TickerSlim is a single ticker_slim subscription update.
//
// Derive's WebSocket `ticker_slim.<inst>.<interval>` channel emits a
// compact wire payload with single-letter fields. The payload is wrapped
// in a `{timestamp, instrument_ticker}` envelope; this type captures that
// envelope and exposes the inner data via the [TickerSlim.Ticker] field.
type TickerSlim struct {
	// Timestamp is the message-emission time in milliseconds.
	Timestamp MillisTime `json:"timestamp"`
	// Ticker is the per-instrument snapshot.
	Ticker InstrumentTickerSlim `json:"instrument_ticker"`
}

// InstrumentTickerSlim is the inner per-instrument payload of a
// `ticker_slim` notification. JSON tags use Derive's compact single-letter
// wire format; Go field names use canonical pascal-case so the type is
// idiomatic to use.
type InstrumentTickerSlim struct {
	// Timestamp is the snapshot's own millisecond timestamp.
	Timestamp MillisTime `json:"t"`

	// BestAskAmount and BestAskPrice are the resting top-ask.
	BestAskAmount Decimal `json:"A"`
	BestAskPrice  Decimal `json:"a"`
	// BestBidAmount and BestBidPrice are the resting top-bid.
	BestBidAmount Decimal `json:"B"`
	BestBidPrice  Decimal `json:"b"`

	// IndexPrice is the underlying oracle price.
	IndexPrice Decimal `json:"I,omitempty"`
	// MarkPrice is the engine-computed mark.
	MarkPrice Decimal `json:"M,omitempty"`

	// FundingRate is the current 1h funding rate (perp instruments only).
	FundingRate Decimal `json:"f,omitempty"`

	// Stats is the rolling 24h volume / OI block. Preserved as raw JSON
	// because Derive's response includes per-margin-type breakdowns whose
	// schema is documented at docs.derive.xyz; decode further if needed.
	Stats json.RawMessage `json:"stats,omitempty"`

	// OptionPricing is the option-specific Greeks/IV block. Preserved as
	// raw JSON because the shape varies by `instrument_type`.
	OptionPricing json.RawMessage `json:"option_pricing,omitempty"`
}
