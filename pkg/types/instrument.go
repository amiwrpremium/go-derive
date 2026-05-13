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

import "github.com/amiwrpremium/go-derive/pkg/enums"

// Instrument describes one tradable contract.
//
// Common fields (Name through IndexPrice) are populated for every instrument
// kind. The kind-specific embedded structs (Perp / Option / ERC20) are
// non-nil only for the matching [enums.InstrumentType]: a perp instrument
// has a non-nil Perp and nil Option/ERC20, and so on. Test the [Type] field
// to know which detail block is populated.
type Instrument struct {
	// Name is the canonical instrument name (e.g. "BTC-PERP",
	// "ETH-25DEC25-65000-C"). Used as a key in pretty much every other
	// API call.
	Name string `json:"instrument_name"`
	// BaseCurrency is the underlying asset symbol (e.g. "BTC", "ETH").
	BaseCurrency string `json:"base_currency"`
	// QuoteCurrency is the asset prices are quoted in — almost always "USDC".
	QuoteCurrency string `json:"quote_currency"`
	// Type identifies which of Perp / Option / ERC20 below is populated.
	Type enums.InstrumentType `json:"instrument_type"`
	// IsActive reports whether the instrument is currently live and tradable.
	IsActive bool `json:"is_active"`
	// TickSize is the minimum price increment.
	TickSize Decimal `json:"tick_size"`
	// MinimumAmount is the smallest order size allowed.
	MinimumAmount Decimal `json:"minimum_amount"`
	// MaximumAmount is the largest order size allowed.
	MaximumAmount Decimal `json:"maximum_amount"`
	// AmountStep is the size increment; sizes must be a whole-number
	// multiple of this value.
	AmountStep Decimal `json:"amount_step"`
	// MarkPrice is the engine's current mark price. Zero until the engine
	// has produced its first mark.
	MarkPrice Decimal `json:"mark_price,omitempty"`
	// IndexPrice is the underlying index price. Zero until the engine has
	// produced its first index print.
	IndexPrice Decimal `json:"index_price,omitempty"`
	// BaseAssetAddress is the on-chain address of the base asset (used by
	// the trade-module signing path).
	BaseAssetAddress Address `json:"base_asset_address,omitempty"`
	// BaseAssetSubID is the per-asset subId Derive uses to disambiguate
	// option strikes / expiries packed into a single ERC-1155.
	BaseAssetSubID string `json:"base_asset_sub_id,omitempty"`

	// MakerFeeRate is the fee rate charged to makers (e.g. "0.0003").
	MakerFeeRate Decimal `json:"maker_fee_rate,omitempty"`
	// TakerFeeRate is the fee rate charged to takers.
	TakerFeeRate Decimal `json:"taker_fee_rate,omitempty"`
	// BaseFee is the flat per-fill fee in quote currency.
	BaseFee Decimal `json:"base_fee,omitempty"`
	// MarkPriceFeeRateCap caps the fee at a fraction of mark price.
	// Nullable on the wire; absent decodes to zero.
	MarkPriceFeeRateCap Decimal `json:"mark_price_fee_rate_cap,omitempty"`

	// ProRataFraction is the fraction of incoming size routed
	// pro-rata to the resting book.
	ProRataFraction Decimal `json:"pro_rata_fraction,omitempty"`
	// ProRataAmountStep is the size increment used by the pro-rata
	// matcher.
	ProRataAmountStep Decimal `json:"pro_rata_amount_step,omitempty"`
	// FIFOMinAllocation is the minimum allocation routed FIFO before
	// pro-rata kicks in.
	FIFOMinAllocation Decimal `json:"fifo_min_allocation,omitempty"`

	// ScheduledActivation is the Unix-seconds time the instrument
	// becomes tradable.
	ScheduledActivation int64 `json:"scheduled_activation,omitempty"`
	// ScheduledDeactivation is the Unix-seconds time the instrument
	// is delisted.
	ScheduledDeactivation int64 `json:"scheduled_deactivation,omitempty"`

	// Perp carries perp-specific fields when [Type] is
	// [enums.InstrumentTypePerp]; nil otherwise.
	Perp *PerpDetails `json:"perp_details,omitempty"`
	// Option carries option-specific fields when [Type] is
	// [enums.InstrumentTypeOption]; nil otherwise.
	Option *OptionDetails `json:"option_details,omitempty"`
	// ERC20 carries ERC-20 spot fields when [Type] is
	// [enums.InstrumentTypeERC20]; nil otherwise.
	ERC20 *ERC20Details `json:"erc20_details,omitempty"`
}

// PerpDetails carries fields specific to perpetual futures contracts.
type PerpDetails struct {
	// IndexName is the index this perp tracks (e.g. "BTC", "ETH").
	IndexName string `json:"index"`
	// MaxLeverage is the maximum leverage allowed for positions on this perp.
	MaxLeverage Decimal `json:"max_leverage,omitempty"`
	// AggregateFundingRate is the cumulative funding rate paid since
	// instrument inception.
	AggregateFundingRate Decimal `json:"aggregate_funding,omitempty"`
	// FundingRate is the most recent per-period funding rate.
	FundingRate Decimal `json:"funding_rate,omitempty"`
	// MaxRatePerHour caps the funding rate at the upper bound per hour.
	MaxRatePerHour Decimal `json:"max_rate_per_hour,omitempty"`
	// MinRatePerHour caps the funding rate at the lower bound per hour.
	MinRatePerHour Decimal `json:"min_rate_per_hour,omitempty"`
	// StaticInterestRate is the engine's static interest-rate input to
	// the funding-rate formula.
	StaticInterestRate Decimal `json:"static_interest_rate,omitempty"`
}

// OptionDetails carries fields specific to options contracts.
type OptionDetails struct {
	// OptionType is call or put.
	OptionType enums.OptionType `json:"option_type"`
	// Strike is the option strike price.
	Strike Decimal `json:"strike"`
	// Expiry is the option expiry timestamp.
	Expiry MillisTime `json:"expiry"`
	// IndexName is the index the option references.
	IndexName string `json:"index"`
	// SettlementPrice is populated after expiry once the option settles.
	SettlementPrice Decimal `json:"settlement_price,omitempty"`
}

// ERC20Details carries fields specific to ERC-20 spot tokens (typically
// collateral assets like USDC, weETH, sUSDe).
type ERC20Details struct {
	// UnderlyingERC20Address is the on-chain address of the wrapped ERC-20.
	UnderlyingERC20Address Address `json:"underlying_erc20_address,omitempty"`
	// Decimals is the underlying ERC-20's decimals (typically 6 for USDC,
	// 18 for weETH).
	Decimals int `json:"decimals,omitempty"`
	// BorrowIndex is the cumulative interest index for borrows.
	BorrowIndex Decimal `json:"borrow_index,omitempty"`
	// SupplyIndex is the cumulative interest index for supplies.
	SupplyIndex Decimal `json:"supply_index,omitempty"`
}
