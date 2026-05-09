// Package derive — domain types for the domain types used in REST and WebSocket
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
package derive

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
	"time"
)

// Address is a 20-byte Ethereum address that JSON-encodes in EIP-55 mixed
// case ("0xAbCd...").
//
// It is a defined type over [common.Address] so callers can convert
// freely with [Address.Common] and the surrounding struct fields stay
// strongly-typed.
//
// The zero value is the all-zero address; use [Address.IsZero] to detect it.
type Address common.Address

// NewAddress parses the hex string s into an [Address]. Both 0x-prefixed and
// unprefixed forms are accepted. The empty string yields the zero address
// with no error so optional fields can decode without ceremony.
func NewAddress(s string) (Address, error) {
	if s == "" {
		return Address{}, nil
	}
	if !common.IsHexAddress(s) {
		return Address{}, fmt.Errorf("types: invalid address %q", s)
	}
	return Address(common.HexToAddress(s)), nil
}

// MustAddress is [NewAddress] that panics on failure. It is appropriate in
// tests and constants where the input is known-good.
func MustAddress(s string) Address {
	a, err := NewAddress(s)
	if err != nil {
		panic(err)
	}
	return a
}

// String returns the EIP-55 mixed-case hex form, including the "0x" prefix.
func (a Address) String() string { return common.Address(a).Hex() }

// Common returns the underlying [common.Address] for interop with
// go-ethereum APIs.
func (a Address) Common() common.Address { return common.Address(a) }

// IsZero reports whether the address equals the zero value (all-zero bytes).
func (a Address) IsZero() bool { return common.Address(a) == (common.Address{}) }

// MarshalJSON encodes the address as a JSON string in EIP-55 form.
func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON decodes a JSON string into an [Address]. The empty string
// yields the zero address; non-string and malformed inputs return an error.
func (a *Address) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !common.IsHexAddress(s) {
		return fmt.Errorf("types: invalid address %q", s)
	}
	*a = Address(common.HexToAddress(s))
	return nil
}

// Collateral is one collateral asset balance for a subaccount.
//
// Each subaccount can hold multiple collaterals; PMRM (portfolio-margin
// risk-managed) subaccounts are restricted to USDC.
type Collateral struct {
	// AssetName is the human-readable symbol (e.g. "USDC", "weETH").
	AssetName string `json:"asset_name"`
	// AssetType identifies the asset class — see [AssetType].
	AssetType AssetType `json:"asset_type"`
	// Currency is the underlying currency (e.g. "USDC", "ETH").
	Currency string `json:"currency,omitempty"`
	// Amount is the balance in the asset's native units.
	Amount Decimal `json:"amount"`
	// MarkPrice is the asset's mark price in quote currency (USDC).
	MarkPrice Decimal `json:"mark_price,omitempty"`
	// MarkValue is the asset balance valued at the current mark.
	MarkValue Decimal `json:"mark_value"`
	// CumulativeInterest is the lifetime interest earned/paid on this asset.
	CumulativeInterest Decimal `json:"cumulative_interest,omitempty"`
	// PendingInterest is interest accrued but not yet settled.
	PendingInterest Decimal `json:"pending_interest,omitempty"`
	// InitialMargin is the asset's contribution to the subaccount's IM.
	InitialMargin Decimal `json:"initial_margin,omitempty"`
	// MaintenanceMargin is the asset's contribution to the subaccount's MM.
	MaintenanceMargin Decimal `json:"maintenance_margin,omitempty"`
}

// Balance summarises a subaccount's value and margin posture in one struct.
//
// SubaccountValue is the headline equity number; InitialMargin and
// MaintenanceMargin set the bands inside which open orders are accepted
// and outside which the engine liquidates.
type Balance struct {
	// SubaccountID identifies the subaccount this balance belongs to.
	SubaccountID int64 `json:"subaccount_id"`
	// SubaccountValue is the total equity (collateral + unrealized PnL +
	// pending funding).
	SubaccountValue Decimal `json:"subaccount_value"`
	// InitialMargin is the margin required to open new orders.
	InitialMargin Decimal `json:"initial_margin"`
	// MaintenanceMargin is the margin floor; breaching it triggers
	// liquidation.
	MaintenanceMargin Decimal `json:"maintenance_margin"`
	// Collaterals is the per-asset balance breakdown.
	Collaterals []Collateral `json:"collaterals"`
	// Positions is the open positions by instrument (omitted by some endpoints).
	Positions []Position `json:"positions,omitempty"`
}

// BalanceUpdate is one entry on the `subaccount.{id}.balances` subscription
// channel. Where Balance is a snapshot, BalanceUpdate is a delta event:
// it carries the [BalanceUpdateType] explaining what caused the
// change (a fill, a deposit, an interest accrual, etc.).
type BalanceUpdate struct {
	// SubaccountID identifies the subaccount this update belongs to.
	SubaccountID int64 `json:"subaccount_id"`
	// AssetName is the affected asset.
	AssetName string `json:"asset_name,omitempty"`
	// AssetType identifies the asset class.
	AssetType AssetType `json:"asset_type,omitempty"`
	// Amount is the new balance after the update.
	Amount Decimal `json:"amount,omitempty"`
	// PreviousAmount is the balance before the update.
	PreviousAmount Decimal `json:"previous_amount,omitempty"`
	// Delta is the signed change.
	Delta Decimal `json:"delta,omitempty"`
	// UpdateType classifies the cause of the update — see
	// [BalanceUpdateType].
	UpdateType BalanceUpdateType `json:"update_type,omitempty"`
	// TxHash is the on-chain transaction hash that generated the update,
	// for update types that involve on-chain settlement.
	TxHash TxHash `json:"tx_hash,omitempty"`
	// TxStatus is the on-chain settlement state.
	TxStatus TxStatus `json:"tx_status,omitempty"`
	// Timestamp is when the update was recorded.
	Timestamp MillisTime `json:"timestamp,omitempty"`
}

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

// Decimal is a fixed-precision decimal number that JSON-encodes as a string.
//
// Derive returns numbers (prices, sizes, fees) as JSON strings to avoid the
// truncation float64 would impose on 18-decimal-place values. [Decimal]
// preserves the round-trip byte-for-byte and supports the full
// shopspring/decimal arithmetic API via [Decimal.Inner].
//
// The zero value is the decimal zero; it is safe to use without
// initialisation.
type Decimal struct{ d decimal.Decimal }

// NewDecimal parses the canonical decimal representation s into a [Decimal].
// Acceptable forms include "0", "1.5", "0.0000000000000000018", "-2.5", and
// scientific notation ("1.5e3"). It returns an error if s is not a valid
// decimal.
func NewDecimal(s string) (Decimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Decimal{}, fmt.Errorf("types: parse decimal %q: %w", s, err)
	}
	return Decimal{d: d}, nil
}

// MustDecimal is [NewDecimal] that panics on failure. It is appropriate in
// constants and tests where the input is known-good and a parse failure
// is a programmer bug.
func MustDecimal(s string) Decimal {
	d, err := NewDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

// DecimalFromInt builds a [Decimal] from a signed integer. It is exact for
// any int64 value.
func DecimalFromInt(n int64) Decimal {
	return Decimal{d: decimal.NewFromInt(n)}
}

// String returns the canonical decimal representation — the same string
// shopspring/decimal would produce, with trailing zeroes stripped.
func (d Decimal) String() string { return d.d.String() }

// Inner returns the underlying shopspring/decimal.Decimal, allowing callers
// to perform arithmetic without a round-trip through string.
//
// The returned value is a copy and is independent of the receiver.
func (d Decimal) Inner() decimal.Decimal { return d.d }

// IsZero reports whether the decimal equals zero.
func (d Decimal) IsZero() bool { return d.d.IsZero() }

// Sign returns -1, 0 or +1 for negative, zero or positive values
// respectively.
func (d Decimal) Sign() int { return d.d.Sign() }

// MarshalJSON encodes the decimal as a JSON string — the form Derive
// expects on the wire.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.d.String())
}

// UnmarshalJSON decodes a JSON string or number into the receiver.
//
// The Derive API always emits strings, but the implementation tolerates
// numeric input for resilience. Empty strings and JSON null leave the
// receiver untouched.
func (d *Decimal) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			return nil
		}
		v, err := decimal.NewFromString(s)
		if err != nil {
			return fmt.Errorf("types: decode decimal %q: %w", s, err)
		}
		d.d = v
		return nil
	}
	v, err := decimal.NewFromString(string(b))
	if err != nil {
		return fmt.Errorf("types: decode decimal %s: %w", b, err)
	}
	d.d = v
	return nil
}

// Instrument describes one tradable contract.
//
// Common fields (Name through IndexPrice) are populated for every instrument
// kind. The kind-specific embedded structs (Perp / Option / ERC20) are
// non-nil only for the matching [InstrumentType]: a perp instrument
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
	Type InstrumentType `json:"instrument_type"`
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

	// Perp carries perp-specific fields when [Type] is
	// [InstrumentTypePerp]; nil otherwise.
	Perp *PerpDetails `json:"perp_details,omitempty"`
	// Option carries option-specific fields when [Type] is
	// [InstrumentTypeOption]; nil otherwise.
	Option *OptionDetails `json:"option_details,omitempty"`
	// ERC20 carries ERC-20 spot fields when [Type] is
	// [InstrumentTypeERC20]; nil otherwise.
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
	AggregateFundingRate Decimal `json:"aggregate_funding_rate,omitempty"`
	// FundingRate is the most recent per-period funding rate.
	FundingRate Decimal `json:"funding_rate,omitempty"`
}

// OptionDetails carries fields specific to options contracts.
type OptionDetails struct {
	// OptionType is call or put.
	OptionType OptionType `json:"option_type"`
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
	// BorrowIndex is the cumulative interest index for borrows.
	BorrowIndex Decimal `json:"borrow_index,omitempty"`
	// SupplyIndex is the cumulative interest index for supplies.
	SupplyIndex Decimal `json:"supply_index,omitempty"`
}

// AuctionBid is one bid placed against a liquidation auction.
type AuctionBid struct {
	// Bidder is the wallet that placed the bid.
	Bidder Address `json:"bidder,omitempty"`
	// Price is the price the bidder offered.
	Price Decimal `json:"price,omitempty"`
	// PercentLiquidated is how much of the position the bid covers.
	PercentLiquidated Decimal `json:"percent_liquidated,omitempty"`
	// Timestamp is when the bid was received.
	Timestamp MillisTime `json:"timestamp,omitempty"`
}

// Liquidation is a liquidation-auction event reported by the engine.
//
// The canonical shape mirrors `derivexyz/cockpit`'s
// `AuctionResultSchema`. Public payloads carry the auction outcome;
// private payloads on the affected subaccount add per-position breakdowns.
type Liquidation struct {
	// AuctionID is the unique server-side auction id.
	AuctionID string `json:"auction_id,omitempty"`
	// AuctionType is "solvent" or "insolvent" — see [AuctionType].
	AuctionType AuctionType `json:"auction_type,omitempty"`
	// SubaccountID is the subaccount being liquidated.
	SubaccountID int64 `json:"subaccount_id,omitempty"`
	// StartTimestamp is when the auction opened.
	StartTimestamp MillisTime `json:"start_timestamp,omitempty"`
	// EndTimestamp is when the auction closed (zero for ongoing).
	EndTimestamp MillisTime `json:"end_timestamp,omitempty"`
	// Bids is the chronological list of bids placed during the auction.
	Bids []AuctionBid `json:"bids,omitempty"`
	// CashReceived is the cash flow into the affected subaccount.
	CashReceived Decimal `json:"cash_received,omitempty"`
	// DiscountPnL is the realized PnL contribution from the auction discount.
	DiscountPnL Decimal `json:"discount_pnl,omitempty"`
	// Fee is the auction fee charged.
	Fee Decimal `json:"fee,omitempty"`
	// PercentLiquidated is how much of the subaccount was liquidated [0, 1].
	PercentLiquidated Decimal `json:"percent_liquidated,omitempty"`
	// RealizedPnL is the total realized PnL across the auction.
	RealizedPnL Decimal `json:"realized_pnl,omitempty"`
	// AmountsLiquidated is per-instrument size that was liquidated.
	AmountsLiquidated map[string]Decimal `json:"amounts_liquidated,omitempty"`
	// PositionsRealizedPnL is per-instrument realized PnL.
	PositionsRealizedPnL map[string]Decimal `json:"positions_realized_pnl,omitempty"`
	// Timestamp is when the engine recorded the event.
	Timestamp MillisTime `json:"timestamp"`
	// TxHash is the on-chain liquidation transaction hash, if available.
	TxHash TxHash `json:"tx_hash,omitempty"`
}

// OrderBookLevel is one [price, amount] pair on either side of an order book.
//
// Derive serializes order-book levels as two-element JSON arrays rather than
// objects, e.g. ["65000", "1.5"]. A custom Marshal/Unmarshal preserves that
// wire format while exposing readable Price / Amount fields at the call site.
type OrderBookLevel struct {
	// Price is the resting limit price.
	Price Decimal
	// Amount is the resting size at that price.
	Amount Decimal
}

// MarshalJSON encodes the level as a [price, amount] JSON array.
func (l OrderBookLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]Decimal{l.Price, l.Amount})
}

// UnmarshalJSON decodes a [price, amount] JSON array into the receiver.
// Object-shaped input returns an error.
func (l *OrderBookLevel) UnmarshalJSON(b []byte) error {
	var arr [2]Decimal
	if err := json.Unmarshal(b, &arr); err != nil {
		return fmt.Errorf("types: orderbook level: %w", err)
	}
	l.Price = arr[0]
	l.Amount = arr[1]
	return nil
}

// OrderBook is a snapshot of the order book at a point in time.
//
// Both sides are sorted: [Bids] is descending by price, [Asks] is ascending.
type OrderBook struct {
	// InstrumentName identifies the market this snapshot belongs to.
	InstrumentName string `json:"instrument_name"`
	// Bids is the buy side of the book.
	Bids []OrderBookLevel `json:"bids"`
	// Asks is the sell side of the book.
	Asks []OrderBookLevel `json:"asks"`
	// Timestamp is the engine-side capture time.
	Timestamp MillisTime `json:"timestamp"`
	// PublishTime is the time the snapshot was published over the wire.
	// Compare against Timestamp to gauge engine-to-client latency.
	PublishTime MillisTime `json:"publish_time,omitempty"`
}

// OrderParams is the request shape for `private/order`.
//
// Most fields map directly to the JSON-RPC schema. The four signing fields
// (Signer, Signature, Nonce, SignatureExpiry) are populated automatically by
// [github.com/amiwrpremium/go-derive.API.PlaceOrder] using
// the configured signer; callers building this struct manually must populate
// them themselves and produce a matching EIP-712 signature.
type OrderParams struct {
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// Direction is buy or sell.
	Direction Direction `json:"direction"`
	// OrderType is limit or market.
	OrderType OrderType `json:"order_type"`
	// TimeInForce is the order's expiry policy.
	TimeInForce TimeInForce `json:"time_in_force,omitempty"`
	// Amount is the order size in base-currency units.
	Amount Decimal `json:"amount"`
	// LimitPrice is the price; for market orders this is the slippage cap.
	LimitPrice Decimal `json:"limit_price"`
	// MaxFee is the maximum acceptable fee paid for this order.
	MaxFee Decimal `json:"max_fee"`
	// SubaccountID is the placing subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Nonce is the order's monotonic anti-replay nonce.
	Nonce uint64 `json:"nonce"`
	// Signer is the signing key's public address (session key, owner, etc.).
	Signer Address `json:"signer"`
	// Signature is the hex-encoded EIP-712 signature over the action data.
	Signature string `json:"signature"`
	// SignatureExpiry is the Unix timestamp (seconds) after which the
	// signature is no longer valid.
	SignatureExpiry int64 `json:"signature_expiry_sec"`

	// Label is a free-form per-order tag, useful for cancel-by-label.
	Label string `json:"label,omitempty"`
	// MMP enrols the order in market-maker protection accounting.
	MMP bool `json:"mmp,omitempty"`
	// ReduceOnly forces the order to only reduce position size, not flip
	// or grow it.
	ReduceOnly bool `json:"reduce_only,omitempty"`
}

// Order is the canonical order record returned by the API. It carries both
// the fields the user supplied and the engine's lifecycle state.
type Order struct {
	// OrderID is the unique server-side id.
	OrderID string `json:"order_id"`
	// SubaccountID is the placing subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// Direction is buy or sell.
	Direction Direction `json:"direction"`
	// OrderType is limit or market.
	OrderType OrderType `json:"order_type"`
	// TimeInForce is the order's expiry policy.
	TimeInForce TimeInForce `json:"time_in_force"`
	// OrderStatus is the current lifecycle state. Once it transitions to
	// a [OrderStatus.Terminal] value, no further updates arrive.
	OrderStatus OrderStatus `json:"order_status"`
	// Amount is the original order size.
	Amount Decimal `json:"amount"`
	// FilledAmount is the cumulative filled size so far.
	FilledAmount Decimal `json:"filled_amount"`
	// LimitPrice is the original limit price.
	LimitPrice Decimal `json:"limit_price"`
	// AveragePrice is the volume-weighted fill price (zero if no fills).
	AveragePrice Decimal `json:"average_price,omitempty"`
	// MaxFee is the original max-fee cap.
	MaxFee Decimal `json:"max_fee"`
	// Nonce is the original nonce.
	Nonce uint64 `json:"nonce"`
	// Signer is the address that signed the order.
	Signer Address `json:"signer"`
	// Label is the user-supplied label (empty if none).
	Label string `json:"label,omitempty"`
	// CancelReason is populated when [OrderStatus] is
	// [OrderStatusCancelled]; empty otherwise.
	CancelReason CancelReason `json:"cancel_reason,omitempty"`
	// MMP indicates the order participated in market-maker-protection
	// accounting.
	MMP bool `json:"mmp,omitempty"`
	// ReduceOnly indicates the order was constrained to reducing position size.
	ReduceOnly bool `json:"reduce_only,omitempty"`
	// IsTransfer indicates the order was a synthetic order created by an
	// internal sub-account transfer rather than a user submission.
	IsTransfer bool `json:"is_transfer,omitempty"`
	// QuoteID links this order to the maker quote it executed against,
	// when the fill came out of an RFQ flow.
	QuoteID string `json:"quote_id,omitempty"`
	// ReplacedOrderID points back to the original order id when this
	// order was created via `private/replace`.
	ReplacedOrderID string `json:"replaced_order_id,omitempty"`
	// OrderFee is the cumulative fee charged on this order's fills.
	OrderFee Decimal `json:"order_fee,omitempty"`
	// SignatureExpiry is the Unix timestamp (seconds) after which the
	// signature is no longer valid.
	SignatureExpiry int64 `json:"signature_expiry_sec,omitempty"`
	// Signature is the EIP-712 action signature attached to the order.
	Signature string `json:"signature,omitempty"`
	// CreationTimestamp is when the engine first saw the order.
	CreationTimestamp MillisTime `json:"creation_timestamp"`
	// LastUpdateTimestamp is the engine's most recent update time.
	LastUpdateTimestamp MillisTime `json:"last_update_timestamp"`
}

// CancelOrderParams identifies an order to cancel via `private/cancel`.
//
// Either OrderID or Label must be set; if both are present the engine
// prefers OrderID. The signing fields are populated automatically when
// using the high-level client.
type CancelOrderParams struct {
	// SubaccountID is the placing subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// InstrumentName scopes the cancel to one market (optional).
	InstrumentName string `json:"instrument_name,omitempty"`
	// OrderID identifies the specific order to cancel.
	OrderID string `json:"order_id,omitempty"`
	// Label cancels every order carrying this label.
	Label string `json:"label,omitempty"`
	// Nonce is the cancel-action nonce.
	Nonce uint64 `json:"nonce,omitempty"`
	// Signer is the signing key's address.
	Signer Address `json:"signer,omitempty"`
	// Signature is the hex EIP-712 signature.
	Signature string `json:"signature,omitempty"`
	// SignatureExpiry is the cancel signature's expiry.
	SignatureExpiry int64 `json:"signature_expiry_sec,omitempty"`
}

// ReplaceOrderParams atomically cancels one order and places another. The
// matching engine guarantees there is no window in which neither order is
// live.
type ReplaceOrderParams struct {
	// OrderIDToCancel is the existing order to drop.
	OrderIDToCancel string `json:"order_id_to_cancel"`
	// NewOrder is the replacement order spec.
	NewOrder OrderParams `json:"new_order"`
}

// ErrInvalidParams is the sentinel returned by every param-DTO Validate
// method when the receiver fails the schema-level checks (required field
// missing, value out of range, enum value not recognised). Wrap with
// errors.Is.
var ErrInvalidParams = errors.New("types: invalid params")

// invalidParam is the package-internal helper for assembling consistent
// validation errors. The returned error wraps [ErrInvalidParams] so
// callers can match without unwrapping each kind separately.
func invalidParam(field, reason string) error {
	return fmt.Errorf("%w: %s: %s", ErrInvalidParams, field, reason)
}

// NewOrderParams constructs an [OrderParams] populated with the four
// always-required fields. Use the With* methods to attach optional
// values; the high-level client supplies the signing fields
// (Signer/Signature/Nonce/SignatureExpiry) at submission time.
//
// The returned struct is unsigned and unvalidated — call [OrderParams.Validate]
// before serialising.
func NewOrderParams(instrument string, side Direction, kind OrderType, amount, limitPrice Decimal) OrderParams {
	return OrderParams{
		InstrumentName: instrument,
		Direction:      side,
		OrderType:      kind,
		Amount:         amount,
		LimitPrice:     limitPrice,
	}
}

// WithInstrument returns a copy with the instrument name set.
func (p OrderParams) WithInstrument(name string) OrderParams { p.InstrumentName = name; return p }

// WithDirection returns a copy with [Direction] set.
func (p OrderParams) WithDirection(d Direction) OrderParams { p.Direction = d; return p }

// WithOrderType returns a copy with [OrderType] set.
func (p OrderParams) WithOrderType(o OrderType) OrderParams { p.OrderType = o; return p }

// WithTimeInForce returns a copy with [TimeInForce] set.
func (p OrderParams) WithTimeInForce(tif TimeInForce) OrderParams {
	p.TimeInForce = tif
	return p
}

// WithAmount returns a copy with the order size set.
func (p OrderParams) WithAmount(amount Decimal) OrderParams { p.Amount = amount; return p }

// WithLimitPrice returns a copy with the limit price set.
func (p OrderParams) WithLimitPrice(price Decimal) OrderParams { p.LimitPrice = price; return p }

// WithMaxFee returns a copy with the max-fee cap set.
func (p OrderParams) WithMaxFee(fee Decimal) OrderParams { p.MaxFee = fee; return p }

// WithSubaccount returns a copy with the subaccount id set.
func (p OrderParams) WithSubaccount(id int64) OrderParams { p.SubaccountID = id; return p }

// WithLabel returns a copy with the per-order label set.
func (p OrderParams) WithLabel(label string) OrderParams { p.Label = label; return p }

// WithMMP returns a copy with market-maker-protection enrolment enabled.
func (p OrderParams) WithMMP() OrderParams { p.MMP = true; return p }

// WithReduceOnly returns a copy with the reduce-only flag enabled.
func (p OrderParams) WithReduceOnly() OrderParams { p.ReduceOnly = true; return p }

// WithSignature returns a copy with the signing-quad set. Callers using
// the high-level client should not call this — the client populates
// these fields itself.
func (p OrderParams) WithSignature(signer Address, signature string, nonce uint64, expiry int64) OrderParams {
	p.Signer = signer
	p.Signature = signature
	p.Nonce = nonce
	p.SignatureExpiry = expiry
	return p
}

// Validate performs schema-level checks on the receiver: required fields
// populated, enum values in range, numeric fields positive. It does not
// validate against an instrument's tick / amount step (those live on
// [Instrument] and require a network round-trip).
//
// Returns nil on success or a wrapped [ErrInvalidParams] describing the
// first failure on a non-nil receiver.
func (p OrderParams) Validate() error {
	if p.InstrumentName == "" {
		return invalidParam("instrument_name", "required")
	}
	if err := p.Direction.Validate(); err != nil {
		return invalidParam("direction", err.Error())
	}
	if err := p.OrderType.Validate(); err != nil {
		return invalidParam("order_type", err.Error())
	}
	if p.TimeInForce != "" {
		if err := p.TimeInForce.Validate(); err != nil {
			return invalidParam("time_in_force", err.Error())
		}
	}
	if p.Amount.Sign() <= 0 {
		return invalidParam("amount", "must be positive")
	}
	if p.LimitPrice.Sign() <= 0 {
		return invalidParam("limit_price", "must be positive")
	}
	if p.MaxFee.Sign() < 0 {
		return invalidParam("max_fee", "must be non-negative")
	}
	if p.SubaccountID < 0 {
		return invalidParam("subaccount_id", "must be non-negative")
	}
	if p.SignatureExpiry < 0 {
		return invalidParam("signature_expiry_sec", "must be non-negative")
	}
	return nil
}

// NewCancelOrderParams constructs a [CancelOrderParams] keyed on either
// an explicit `OrderID` (recommended) or a `Label` (cancels every order
// with that label). The signing fields are filled in by the high-level
// client at submit time.
func NewCancelOrderParams(subaccountID int64) CancelOrderParams {
	return CancelOrderParams{SubaccountID: subaccountID}
}

// WithOrderID returns a copy targeting one specific order id.
func (p CancelOrderParams) WithOrderID(id string) CancelOrderParams { p.OrderID = id; return p }

// WithLabel returns a copy targeting every order carrying the label.
func (p CancelOrderParams) WithLabel(label string) CancelOrderParams { p.Label = label; return p }

// WithInstrument returns a copy scoping the cancel to one market.
func (p CancelOrderParams) WithInstrument(name string) CancelOrderParams {
	p.InstrumentName = name
	return p
}

// WithSignature returns a copy with the signing-quad set.
func (p CancelOrderParams) WithSignature(signer Address, signature string, nonce uint64, expiry int64) CancelOrderParams {
	p.Signer = signer
	p.Signature = signature
	p.Nonce = nonce
	p.SignatureExpiry = expiry
	return p
}

// Validate performs schema-level checks on the receiver. Either OrderID
// or Label must be set; SubaccountID must be non-negative.
func (p CancelOrderParams) Validate() error {
	if p.SubaccountID < 0 {
		return invalidParam("subaccount_id", "must be non-negative")
	}
	if p.OrderID == "" && p.Label == "" {
		return invalidParam("order_id|label", "one of order_id or label is required")
	}
	if p.SignatureExpiry < 0 {
		return invalidParam("signature_expiry_sec", "must be non-negative")
	}
	return nil
}

// NewReplaceOrderParams constructs a [ReplaceOrderParams] for an atomic
// cancel-and-place. orderIDToCancel is required; the new order spec is
// built up via the OrderParams builder.
func NewReplaceOrderParams(orderIDToCancel string, newOrder OrderParams) ReplaceOrderParams {
	return ReplaceOrderParams{
		OrderIDToCancel: orderIDToCancel,
		NewOrder:        newOrder,
	}
}

// WithOrderIDToCancel returns a copy with the cancel target set.
func (p ReplaceOrderParams) WithOrderIDToCancel(id string) ReplaceOrderParams {
	p.OrderIDToCancel = id
	return p
}

// WithNewOrder returns a copy with the replacement order spec set.
func (p ReplaceOrderParams) WithNewOrder(o OrderParams) ReplaceOrderParams {
	p.NewOrder = o
	return p
}

// Validate performs schema-level checks on the receiver: cancel-target
// must be present, replacement-order spec must validate.
func (p ReplaceOrderParams) Validate() error {
	if p.OrderIDToCancel == "" {
		return invalidParam("order_id_to_cancel", "required")
	}
	if err := p.NewOrder.Validate(); err != nil {
		return fmt.Errorf("new_order: %w", err)
	}
	return nil
}

// Page wraps the server-side pagination shape used by every Derive list
// endpoint. The fields mirror the JSON response exactly — Derive returns
// just the totals and lets the caller track which page they asked for.
type Page struct {
	// NumPages is the total number of pages available for the query.
	NumPages int `json:"num_pages"`
	// Count is the total number of records across all pages.
	Count int `json:"count"`
}

// PageRequest is the common pagination input.
//
// Both Page and PageSize are 1-indexed; zero values are omitted on the wire
// so the server's defaults apply.
type PageRequest struct {
	// Page selects which page to fetch (1-indexed). Zero asks for the default.
	Page int `json:"page,omitempty"`
	// PageSize sets how many records per page. Zero asks for the default.
	PageSize int `json:"page_size,omitempty"`
}

// NewPageRequest constructs a [PageRequest] with both fields zero, which
// asks the server for its defaults.
func NewPageRequest() PageRequest { return PageRequest{} }

// WithPage returns a copy with the 1-indexed page set.
func (p PageRequest) WithPage(page int) PageRequest { p.Page = page; return p }

// WithPageSize returns a copy with the page size set.
func (p PageRequest) WithPageSize(size int) PageRequest { p.PageSize = size; return p }

// Validate enforces the schema: both fields must be non-negative.
// A zero in either slot is interpreted as "use the server default" by
// the json `omitempty` tag.
func (p PageRequest) Validate() error {
	if p.Page < 0 {
		return invalidParam("page", "must be non-negative")
	}
	if p.PageSize < 0 {
		return invalidParam("page_size", "must be non-negative")
	}
	return nil
}

// Position is a held position in one instrument on a subaccount.
//
// Amount is signed: positive for long, negative for short, zero for flat.
// Most numeric fields are denominated in the quote currency (USDC for the
// vast majority of Derive markets); Amount itself is in base-currency units.
//
// Greeks (Delta/Gamma/Theta/Vega) are populated for option positions and
// zero for perp/erc20 positions.
type Position struct {
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// InstrumentType identifies whether this is a perp, option or ERC20.
	InstrumentType InstrumentType `json:"instrument_type"`
	// CreationTimestamp is when the position first appeared on the engine.
	CreationTimestamp MillisTime `json:"creation_timestamp,omitempty"`
	// Amount is the signed position size (positive=long, negative=short).
	Amount Decimal `json:"amount"`
	// AveragePrice is the volume-weighted entry price.
	AveragePrice Decimal `json:"average_price"`
	// MarkPrice is the engine's current mark.
	MarkPrice Decimal `json:"mark_price"`
	// MarkValue is the position's mark-to-market value in quote currency.
	MarkValue Decimal `json:"mark_value"`
	// IndexPrice is the underlying index price (zero if not yet computed).
	IndexPrice Decimal `json:"index_price,omitempty"`
	// Leverage is the position's effective leverage.
	Leverage Decimal `json:"leverage,omitempty"`
	// LiquidationPrice is the price at which the engine would liquidate
	// (zero if no liquidation risk).
	LiquidationPrice Decimal `json:"liquidation_price,omitempty"`

	// InitialMargin is the engine's initial-margin requirement for this
	// position alone.
	InitialMargin Decimal `json:"initial_margin,omitempty"`
	// MaintenanceMargin is the maintenance-margin requirement.
	MaintenanceMargin Decimal `json:"maintenance_margin,omitempty"`
	// OpenOrdersMargin is the margin reserved for open orders against this position.
	OpenOrdersMargin Decimal `json:"open_orders_margin,omitempty"`

	// UnrealizedPNL is the mark-to-market PnL.
	UnrealizedPNL Decimal `json:"unrealized_pnl"`
	// RealizedPNL is the cumulative realized PnL across closes.
	RealizedPNL Decimal `json:"realized_pnl"`

	// CumulativeFunding is the total funding paid/received over the
	// position's lifetime (perps only).
	CumulativeFunding Decimal `json:"cumulative_funding,omitempty"`
	// PendingFunding is funding accrued since the last settlement.
	PendingFunding Decimal `json:"pending_funding,omitempty"`
	// NetSettlements is the cumulative net of perp / option settlements
	// applied to the position.
	NetSettlements Decimal `json:"net_settlements,omitempty"`

	// Delta, Gamma, Theta, Vega are the option greeks (option positions
	// only; zero for perp / erc20).
	Delta Decimal `json:"delta,omitempty"`
	Gamma Decimal `json:"gamma,omitempty"`
	Theta Decimal `json:"theta,omitempty"`
	Vega  Decimal `json:"vega,omitempty"`
}

// RFQLeg is one leg of a multi-leg RFQ.
//
// Multi-leg RFQs are how Derive supports option spreads, calendars, etc.
// Each leg references its own instrument, direction and amount; legs must
// be unique by instrument (see [github.com/amiwrpremium/go-derive/pkg/errors.CodeLegInstrumentsNotUnique]).
type RFQLeg struct {
	// InstrumentName identifies the leg's market.
	InstrumentName string `json:"instrument_name"`
	// Direction is buy or sell on this leg.
	Direction Direction `json:"direction"`
	// Amount is the leg's size in base-currency units.
	Amount Decimal `json:"amount"`
}

// Validate performs schema-level checks on the receiver: instrument
// non-empty, direction in range, amount positive. Returns nil on
// success or a wrapped [ErrInvalidParams].
func (l RFQLeg) Validate() error {
	if l.InstrumentName == "" {
		return invalidParam("instrument_name", "required")
	}
	if err := l.Direction.Validate(); err != nil {
		return invalidParam("direction", err.Error())
	}
	if l.Amount.Sign() <= 0 {
		return invalidParam("amount", "must be positive")
	}
	return nil
}

// QuoteLeg is a priced leg attached to a maker's [Quote] response.
//
// Distinct from [RFQLeg] because RFQs don't carry per-leg prices —
// quotes do.
type QuoteLeg struct {
	// InstrumentName identifies the leg's market.
	InstrumentName string `json:"instrument_name"`
	// Direction is the maker's side on this leg.
	Direction Direction `json:"direction"`
	// Amount is the leg's size in base-currency units.
	Amount Decimal `json:"amount"`
	// Price is the per-leg price the maker is committing to.
	Price Decimal `json:"price"`
}

// RFQ is a Request-For-Quote initiated by a taker.
//
// The taker broadcasts the RFQ to whitelisted makers; makers respond with
// [Quote] objects. The taker selects a quote and executes via
// `private/execute_quote`.
type RFQ struct {
	// RFQID is the unique server-side id.
	RFQID string `json:"rfq_id"`
	// SubaccountID is the taker's subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Status is the current lifecycle state.
	Status QuoteStatus `json:"status"`
	// CancelReason is set when Status is QuoteStatusCancelled.
	CancelReason CancelReason `json:"cancel_reason,omitempty"`
	// Legs is the per-instrument breakdown (no per-leg prices on RFQs).
	Legs []RFQLeg `json:"legs"`
	// MaxFee is the cap on total fee the taker is willing to pay.
	MaxFee Decimal `json:"max_total_fee,omitempty"`
	// MinPrice and MaxPrice constrain the price band the taker accepts.
	MinPrice Decimal `json:"min_price,omitempty"`
	MaxPrice Decimal `json:"max_price,omitempty"`
	// CreationTimestamp is when the RFQ was first received.
	CreationTimestamp MillisTime `json:"creation_timestamp"`
	// LastUpdateTimestamp is the most recent state change.
	LastUpdateTimestamp MillisTime `json:"last_update_timestamp"`
}

// Quote is a market-maker response to an [RFQ].
//
// The full canonical shape mirrors `derivexyz/cockpit`'s
// `QuoteResultSchema` — every field below appears on the wire.
type Quote struct {
	// QuoteID is the unique server-side id.
	QuoteID string `json:"quote_id"`
	// RFQID identifies the RFQ this quote responds to.
	RFQID string `json:"rfq_id"`
	// SubaccountID is the maker's subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Direction is the maker's side; the taker's fill is the opposite.
	Direction Direction `json:"direction"`
	// Legs is the per-instrument breakdown (must match the RFQ legs).
	Legs []QuoteLeg `json:"legs"`
	// LegsHash is a server-side hash that ties the quote to a specific
	// leg ordering — useful for replay protection.
	LegsHash string `json:"legs_hash,omitempty"`
	// Price is the all-in net price for the package.
	Price Decimal `json:"price,omitempty"`
	// Fee is the fee charged on the quote (when filled).
	Fee Decimal `json:"fee,omitempty"`
	// MaxFee is the maker's cap on the per-fill fee.
	MaxFee Decimal `json:"max_fee,omitempty"`
	// LiquidityRole identifies the quote as maker-side (always "maker").
	LiquidityRole LiquidityRole `json:"liquidity_role,omitempty"`
	// Status is the current lifecycle state.
	Status QuoteStatus `json:"status"`
	// CancelReason is set when Status is QuoteStatusCancelled.
	CancelReason CancelReason `json:"cancel_reason,omitempty"`
	// MMP indicates the quote participated in market-maker-protection
	// accounting on the maker's subaccount.
	MMP bool `json:"mmp,omitempty"`
	// Label is the maker's free-form per-quote tag.
	Label string `json:"label,omitempty"`
	// Nonce is the maker's signed nonce.
	Nonce uint64 `json:"nonce,omitempty"`
	// Signer is the address that signed the quote action.
	Signer Address `json:"signer,omitempty"`
	// Signature is the EIP-712 action signature.
	Signature string `json:"signature,omitempty"`
	// SignatureExpiry is the Unix timestamp (seconds) past which the
	// signature is rejected.
	SignatureExpiry int64 `json:"signature_expiry_sec,omitempty"`
	// TxHash is the on-chain settlement transaction hash, set after the
	// quote is executed.
	TxHash TxHash `json:"tx_hash,omitempty"`
	// TxStatus is the on-chain settlement state.
	TxStatus TxStatus `json:"tx_status,omitempty"`
	// CreationTimestamp is when the quote was received.
	CreationTimestamp MillisTime `json:"creation_timestamp"`
	// LastUpdateTimestamp is the most recent state change.
	LastUpdateTimestamp MillisTime `json:"last_update_timestamp,omitempty"`
}

// SpotFeed is one update from the `spot_feed.{currency}` WebSocket channel.
//
// Derive's oracle feed delivers a per-currency snapshot of mark and 24h-prior
// prices. Use it for liquidation monitoring, basis calculations, or any
// risk surface that needs an oracle anchor independent of the order book.
type SpotFeed struct {
	// Timestamp is the message-emission time in milliseconds.
	Timestamp MillisTime `json:"timestamp"`
	// Feeds is keyed by currency symbol (e.g. "BTC", "ETH") and carries
	// one [SpotFeedEntry] per currency reported in the message. With the
	// per-currency subscription pattern there is usually exactly one entry.
	Feeds map[string]SpotFeedEntry `json:"feeds"`
}

// SpotFeedEntry is the per-currency oracle reading inside a [SpotFeed].
type SpotFeedEntry struct {
	// Price is the current oracle price (decimal string in quote units).
	Price Decimal `json:"price"`
	// Confidence is the oracle confidence score (decimal string in [0, 1]).
	Confidence Decimal `json:"confidence"`
	// PricePrevDaily is the price 24 hours prior, used to derive the
	// 24h delta on the UI without an extra round trip.
	PricePrevDaily Decimal `json:"price_prev_daily"`
	// ConfidencePrevDaily is the oracle confidence 24 hours prior.
	ConfidencePrevDaily Decimal `json:"confidence_prev_daily"`
	// TimestampPrevDaily is the millisecond timestamp of the 24h-prior
	// reading.
	TimestampPrevDaily MillisTime `json:"timestamp_prev_daily"`
}

// SubAccount is a snapshot of one subaccount as returned by GetSubaccount.
//
// A wallet has one or more subaccounts, each with its own positions,
// collateral and margin state. Subaccounts isolate risk: a liquidation in
// one does not cascade into another.
type SubAccount struct {
	// SubaccountID is the unique numeric id.
	SubaccountID int64 `json:"subaccount_id"`
	// OwnerAddress is the smart-account owner that controls this subaccount.
	OwnerAddress Address `json:"owner_address"`
	// MarginType is "PM" (portfolio margin), "SM" (standard margin), etc.
	MarginType string `json:"margin_type"`
	// IsUnderLiquidation is true when the engine is actively liquidating
	// the subaccount.
	IsUnderLiquidation bool `json:"is_under_liquidation"`
	// SubaccountValue is the total equity.
	SubaccountValue Decimal `json:"subaccount_value"`
	// InitialMargin is the margin required to open new orders.
	InitialMargin Decimal `json:"initial_margin"`
	// MaintenanceMargin is the liquidation floor.
	MaintenanceMargin Decimal `json:"maintenance_margin"`
	// OpenOrders is the list of currently-open orders.
	OpenOrders []Order `json:"open_orders,omitempty"`
	// Positions is the list of open positions.
	Positions []Position `json:"positions,omitempty"`
	// Collaterals is the per-asset collateral breakdown.
	Collaterals []Collateral `json:"collaterals,omitempty"`
}

// Ticker is the public market summary for one instrument: top-of-book, marks,
// and depth at 5%.
type Ticker struct {
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// InstrumentType is "perp", "option", or "erc20".
	InstrumentType string `json:"instrument_type,omitempty"`
	// IsActive reports whether the instrument is currently open for trading.
	IsActive bool `json:"is_active,omitempty"`

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

	// Timestamp is when this ticker snapshot was produced.
	Timestamp MillisTime `json:"timestamp"`
}

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
	// schema is documented at docs.xyz; decode further if needed.
	Stats json.RawMessage `json:"stats,omitempty"`

	// OptionPricing is the option-specific Greeks/IV block. Preserved as
	// raw JSON because the shape varies by `instrument_type`.
	OptionPricing json.RawMessage `json:"option_pricing,omitempty"`
}

// MillisTime is a time.Time that round-trips as integer milliseconds since
// the Unix epoch — Derive's preferred timestamp format on every JSON-RPC
// payload.
//
// The zero value is the zero time.Time; use [time.Time.IsZero] on the
// underlying [MillisTime.Time] to detect it.
type MillisTime struct {
	// T is the underlying time. Use [MillisTime.Time] in callers; this
	// field is exported only so that struct literals are convenient.
	T time.Time
}

// NewMillisTime wraps a [time.Time] as a [MillisTime].
func NewMillisTime(t time.Time) MillisTime { return MillisTime{T: t} }

// Time returns the underlying [time.Time].
func (m MillisTime) Time() time.Time { return m.T }

// Millis returns the time as milliseconds since the Unix epoch.
func (m MillisTime) Millis() int64 { return m.T.UnixMilli() }

// MarshalJSON encodes the time as an integer count of milliseconds since
// the Unix epoch.
func (m MillisTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.T.UnixMilli())
}

// UnmarshalJSON decodes either a JSON number or a JSON string of integer
// milliseconds. Empty strings and JSON null leave the receiver as the zero
// value.
func (m *MillisTime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			return nil
		}
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		m.T = time.UnixMilli(n)
		return nil
	}
	var n int64
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	m.T = time.UnixMilli(n)
	return nil
}

// Trade is one filled execution.
//
// Public-trade payloads (from `public/get_trade_history` and the
// `trades.{instrument}` channel) populate the public fields. Private
// trade payloads additionally populate OrderID, SubaccountID, Fee,
// LiquidityRole, RealizedPnL, plus the on-chain settlement fields
// (TxStatus, TxHash, Wallet) and the optional QuoteID / RFQID linking
// when the fill came out of an RFQ flow.
type Trade struct {
	// TradeID is the unique server-side trade id.
	TradeID string `json:"trade_id"`
	// OrderID identifies the order on the user's side of the fill (private
	// only; empty for public-trade payloads).
	OrderID string `json:"order_id,omitempty"`
	// SubaccountID identifies the user's subaccount (private only).
	SubaccountID int64 `json:"subaccount_id,omitempty"`
	// Wallet is the signer wallet that owns the trade. Public-trade payloads
	// include this for tape attribution.
	Wallet Address `json:"wallet,omitempty"`
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// Direction is the user's side: buy means the user took the ask;
	// sell means the user hit the bid.
	Direction Direction `json:"direction"`
	// TradePrice is the printed fill price.
	TradePrice Decimal `json:"trade_price"`
	// TradeAmount is the filled size in base-currency units.
	TradeAmount Decimal `json:"trade_amount"`
	// MarkPrice is the engine's mark price at the time of the fill.
	MarkPrice Decimal `json:"mark_price"`
	// IndexPrice is the underlying index price at the time of the fill.
	IndexPrice Decimal `json:"index_price,omitempty"`
	// Fee is the fee paid (private only).
	Fee Decimal `json:"trade_fee,omitempty"`
	// ExpectedRebate is the maker rebate the engine projects for this fill,
	// if any. Public payloads carry this for tape transparency.
	ExpectedRebate Decimal `json:"expected_rebate,omitempty"`
	// ExtraFee is any additional fee applied (e.g. settlement).
	ExtraFee Decimal `json:"extra_fee,omitempty"`
	// LiquidityRole identifies whether the user was maker or taker (private only).
	LiquidityRole LiquidityRole `json:"liquidity_role,omitempty"`
	// RealizedPnL is the realized PnL from this fill (private only).
	RealizedPnL Decimal `json:"realized_pnl,omitempty"`
	// RealizedPnLExclFees is realized PnL with fees excluded (private only).
	RealizedPnLExclFees Decimal `json:"realized_pnl_excl_fees,omitempty"`
	// QuoteID links the trade to the quote it executed against, when the
	// fill came out of an RFQ flow.
	QuoteID string `json:"quote_id,omitempty"`
	// RFQID links the trade back to the originating RFQ, when applicable.
	RFQID string `json:"rfq_id,omitempty"`
	// TxHash is the on-chain transaction hash that carried the trade
	// settlement (populated once the engine submits the tx).
	TxHash TxHash `json:"tx_hash,omitempty"`
	// TxStatus reports the on-chain settlement state — see
	// [TxStatus] for the allowed values.
	TxStatus TxStatus `json:"tx_status,omitempty"`
	// Timestamp is the engine-side execution time.
	Timestamp MillisTime `json:"timestamp"`
}

// DepositTx records a single deposit into a subaccount.
//
// Returned by `private/get_deposit_history`; also delivered on the
// account-balance channel as deposits finalize.
type DepositTx struct {
	// TxHash is the on-chain deposit transaction hash.
	TxHash TxHash `json:"tx_hash"`
	// Asset is the deposited asset's symbol (e.g. "USDC").
	Asset string `json:"asset"`
	// Amount is the deposited quantity.
	Amount Decimal `json:"amount"`
	// SubaccountID is the credited subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Status is the lifecycle state ("pending", "completed", etc.).
	Status string `json:"status"`
	// Timestamp is when the deposit was first observed.
	Timestamp MillisTime `json:"timestamp"`
}

// WithdrawTx records a single withdrawal from a subaccount.
//
// Withdrawals are two-phase: first the subaccount is debited (status
// pending), then the on-chain transfer is dispatched (status completed).
type WithdrawTx struct {
	// TxHash is the on-chain withdrawal transaction hash.
	TxHash TxHash `json:"tx_hash"`
	// Asset is the withdrawn asset's symbol.
	Asset string `json:"asset"`
	// Amount is the withdrawn quantity.
	Amount Decimal `json:"amount"`
	// SubaccountID is the debited subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Status is the lifecycle state.
	Status string `json:"status"`
	// Timestamp is when the withdrawal was first observed.
	Timestamp MillisTime `json:"timestamp"`
}

// TxHash is a 32-byte transaction hash that JSON-encodes as 0x-prefixed hex.
// It is used for deposit/withdraw acknowledgements and liquidation events.
type TxHash common.Hash

// NewTxHash parses a 0x-prefixed 66-character hex string into a [TxHash].
// The empty string yields the zero hash.
func NewTxHash(s string) (TxHash, error) {
	if s == "" {
		return TxHash{}, nil
	}
	if !strings.HasPrefix(s, "0x") || len(s) != 66 {
		return TxHash{}, fmt.Errorf("types: invalid tx hash %q", s)
	}
	return TxHash(common.HexToHash(s)), nil
}

// String returns the 0x-prefixed lowercase-hex representation.
func (h TxHash) String() string { return common.Hash(h).Hex() }

// IsZero reports whether the hash is all zeros.
func (h TxHash) IsZero() bool { return common.Hash(h) == (common.Hash{}) }

// MarshalJSON encodes the hash as a JSON string.
func (h TxHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON decodes a JSON string into a [TxHash]. The empty string
// yields the zero hash; malformed input returns an error.
func (h *TxHash) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	v, err := NewTxHash(s)
	if err != nil {
		return err
	}
	*h = v
	return nil
}
