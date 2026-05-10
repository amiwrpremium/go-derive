// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the per-event payload of Derive's full `ticker`
// WebSocket channel — a fat snapshot that combines every instrument-
// metadata field with the live top-of-book / mark / index numbers.
//
// The compact wire variant (`ticker_slim`) lives in [TickerSlim] /
// [InstrumentTickerSlim].
package types

import (
	"encoding/json"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

// InstrumentTickerFeed is the envelope emitted by the `ticker`
// WebSocket channel: an [InstrumentTicker] snapshot plus the
// ticker's emission timestamp.
//
// Mirrors the publisher-data envelope documented at
// docs.derive.xyz/reference/ticker-instrument_name-interval.
type InstrumentTickerFeed struct {
	// Timestamp is the emission timestamp (millisecond Unix epoch).
	Timestamp MillisTime `json:"timestamp"`
	// Ticker is the full per-instrument snapshot.
	Ticker InstrumentTicker `json:"instrument_ticker"`
}

// InstrumentTicker is the full ticker payload emitted by the
// `ticker.{instrument}.{interval}` channel. It folds together the
// instrument's metadata (sizes, fees, schedule) and live market
// data (top-of-book, marks, index, price bands) in a single struct.
//
// Mirrors the per-instrument ticker payload documented at
// docs.derive.xyz/reference/ticker-instrument_name-interval.
//
// The four nested complex blocks are kept as [json.RawMessage] for
// the same reasons [InstrumentTickerSlim.Stats] /
// [InstrumentTickerSlim.OptionPricing] are: their schemas are rich
// enough that decoding them in-band would expand this type's surface
// well beyond the SDK's idiom — decode further at the call site if
// needed.
type InstrumentTicker struct {
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// InstrumentType is "perp", "option", or "erc20".
	InstrumentType enums.InstrumentType `json:"instrument_type"`
	// IsActive reports whether the instrument is currently open for
	// trading within its activation window.
	IsActive bool `json:"is_active"`

	// BaseCurrency is the underlying asset symbol (e.g. "BTC", "ETH").
	BaseCurrency string `json:"base_currency"`
	// QuoteCurrency is the asset prices are quoted in (typically
	// "USDC" for options, "USD" for perps).
	QuoteCurrency string `json:"quote_currency"`

	// BaseAssetAddress is the on-chain address of the base asset
	// (used by the trade-module signing path).
	BaseAssetAddress Address `json:"base_asset_address"`
	// BaseAssetSubID is the per-asset subId Derive uses to
	// disambiguate option strikes / expiries packed into a single
	// ERC-1155.
	BaseAssetSubID string `json:"base_asset_sub_id"`

	// AmountStep is the size increment.
	AmountStep Decimal `json:"amount_step"`
	// MinimumAmount is the smallest order size allowed.
	MinimumAmount Decimal `json:"minimum_amount"`
	// MaximumAmount is the largest order size allowed.
	MaximumAmount Decimal `json:"maximum_amount"`
	// TickSize is the minimum price increment.
	TickSize Decimal `json:"tick_size"`

	// BaseFee is the dollar base fee added to every taker order.
	BaseFee Decimal `json:"base_fee"`
	// MakerFeeRate is the percent-of-spot maker fee rate.
	MakerFeeRate Decimal `json:"maker_fee_rate"`
	// TakerFeeRate is the percent-of-spot taker fee rate.
	TakerFeeRate Decimal `json:"taker_fee_rate"`
	// MarkPriceFeeRateCap is the option-price fee cap (e.g. 12.5 %).
	// Nullable on the wire — zero when not applicable.
	MarkPriceFeeRateCap Decimal `json:"mark_price_fee_rate_cap,omitempty"`

	// ScheduledActivation is the timestamp at which the instrument
	// becomes / became active (Unix seconds).
	ScheduledActivation int64 `json:"scheduled_activation"`
	// ScheduledDeactivation is the planned deactivation time, or far
	// in the future for evergreen instruments.
	ScheduledDeactivation int64 `json:"scheduled_deactivation"`

	// BestBidPrice is the highest resting bid.
	BestBidPrice Decimal `json:"best_bid_price"`
	// BestBidAmount is the size resting at [BestBidPrice].
	BestBidAmount Decimal `json:"best_bid_amount"`
	// BestAskPrice is the lowest resting ask.
	BestAskPrice Decimal `json:"best_ask_price"`
	// BestAskAmount is the size resting at [BestAskPrice].
	BestAskAmount Decimal `json:"best_ask_amount"`

	// MarkPrice is the engine's mark price.
	MarkPrice Decimal `json:"mark_price"`
	// IndexPrice is the underlying index price.
	IndexPrice Decimal `json:"index_price"`
	// MinPrice is the engine-enforced lower price band.
	MinPrice Decimal `json:"min_price"`
	// MaxPrice is the engine-enforced upper price band.
	MaxPrice Decimal `json:"max_price"`

	// Timestamp is the engine's snapshot timestamp (millisecond
	// Unix epoch).
	Timestamp MillisTime `json:"timestamp"`

	// OptionDetails is the per-option metadata block. Nullable —
	// empty for non-option instruments. Wire shape mirrors
	// `OptionPublicDetailsSchema`.
	OptionDetails json.RawMessage `json:"option_details,omitempty"`
	// PerpDetails is the per-perp metadata block. Nullable — empty
	// for non-perp instruments.
	PerpDetails json.RawMessage `json:"perp_details,omitempty"`
	// OptionPricing is the option-pricing block (greeks, IV).
	// Nullable — empty for non-option instruments.
	OptionPricing json.RawMessage `json:"option_pricing,omitempty"`
	// Stats is the rolling 24h trading stats (volume, OI,
	// percent_change). Wire shape mirrors
	// `AggregateTradingStatsSchema`.
	Stats json.RawMessage `json:"stats"`
}
