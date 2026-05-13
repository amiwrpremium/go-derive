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
//
// Wire shape: a flat object combining the instrument's static metadata
// (fees, schedule, asset addresses) with the live market state (top-of-book,
// marks, OI). The kind-specific blocks (Perp/Option/ERC20) at the bottom are
// non-nil only for the matching [InstrumentType].
type Ticker struct {
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// InstrumentType is "perp", "option", or "erc20".
	InstrumentType string `json:"instrument_type,omitempty"`
	// IsActive reports whether the instrument is currently open for trading.
	IsActive bool `json:"is_active,omitempty"`

	// BaseCurrency is the underlying asset symbol (e.g. "BTC", "ETH").
	BaseCurrency string `json:"base_currency,omitempty"`
	// QuoteCurrency is the asset prices are quoted in — almost always "USDC".
	QuoteCurrency string `json:"quote_currency,omitempty"`
	// BaseAssetAddress is the on-chain address of the base asset.
	BaseAssetAddress Address `json:"base_asset_address,omitempty"`
	// BaseAssetSubID is the per-asset subId.
	BaseAssetSubID string `json:"base_asset_sub_id,omitempty"`

	// TickSize is the minimum price increment.
	TickSize Decimal `json:"tick_size,omitempty"`
	// AmountStep is the size increment.
	AmountStep Decimal `json:"amount_step,omitempty"`
	// MinimumAmount is the smallest order size allowed.
	MinimumAmount Decimal `json:"minimum_amount,omitempty"`
	// MaximumAmount is the largest order size allowed.
	MaximumAmount Decimal `json:"maximum_amount,omitempty"`

	// MakerFeeRate is the fee rate charged to makers (e.g. "0.0003").
	MakerFeeRate Decimal `json:"maker_fee_rate,omitempty"`
	// TakerFeeRate is the fee rate charged to takers.
	TakerFeeRate Decimal `json:"taker_fee_rate,omitempty"`
	// BaseFee is the flat per-fill fee in quote currency.
	BaseFee Decimal `json:"base_fee,omitempty"`
	// MarkPriceFeeRateCap caps the fee at a fraction of mark price.
	MarkPriceFeeRateCap Decimal `json:"mark_price_fee_rate_cap,omitempty"`

	// ProRataFraction is the fraction of incoming size routed pro-rata.
	ProRataFraction Decimal `json:"pro_rata_fraction,omitempty"`
	// ProRataAmountStep is the size increment used by the pro-rata matcher.
	ProRataAmountStep Decimal `json:"pro_rata_amount_step,omitempty"`
	// FIFOMinAllocation is the minimum allocation routed FIFO before
	// pro-rata kicks in.
	FIFOMinAllocation Decimal `json:"fifo_min_allocation,omitempty"`

	// ScheduledActivation is the Unix-seconds activation time.
	ScheduledActivation int64 `json:"scheduled_activation,omitempty"`
	// ScheduledDeactivation is the Unix-seconds delisting time.
	ScheduledDeactivation int64 `json:"scheduled_deactivation,omitempty"`

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

	// Stats is the rolling 24h volume / OI / price-change block. Preserved
	// as raw JSON because the shape includes per-margin-type breakdowns.
	Stats json.RawMessage `json:"stats,omitempty"`

	// OptionPricing carries IV / greeks for option tickers. Nullable for
	// non-option instruments.
	OptionPricing json.RawMessage `json:"option_pricing,omitempty"`

	// Perp carries perp-specific fields when InstrumentType is "perp".
	Perp *PerpDetails `json:"perp_details,omitempty"`
	// Option carries option-specific fields when InstrumentType is "option".
	Option *OptionDetails `json:"option_details,omitempty"`
	// ERC20 carries ERC-20 spot fields when InstrumentType is "erc20".
	ERC20 *ERC20Details `json:"erc20_details,omitempty"`

	// Timestamp is when this ticker snapshot was produced.
	Timestamp MillisTime `json:"timestamp"`
}
