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
	Direction enums.Direction `json:"direction"`
	// OrderType is limit or market.
	OrderType enums.OrderType `json:"order_type"`
	// TimeInForce is the order's expiry policy.
	TimeInForce enums.TimeInForce `json:"time_in_force"`
	// OrderStatus is the current lifecycle state. Once it transitions to
	// a [enums.OrderStatus.Terminal] value, no further updates arrive.
	OrderStatus enums.OrderStatus `json:"order_status"`
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
	// [enums.OrderStatusCancelled]; empty otherwise.
	CancelReason enums.CancelReason `json:"cancel_reason,omitempty"`
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
	// order was created via `private/replace`. Nullable on the wire.
	ReplacedOrderID string `json:"replaced_order_id,omitempty"`
	// OrderFee is the cumulative fee charged on this order's fills.
	OrderFee Decimal `json:"order_fee,omitempty"`
	// ExtraFee is any extra USDC fee added by a referring client
	// (included in [OrderFee]). Nullable on the wire — zero-value
	// [Decimal] when absent.
	ExtraFee Decimal `json:"extra_fee,omitempty"`
	// SignedLimitPrice is the wire-side `signed_limit_price` —
	// nullable; carries the limit price the maker actually signed
	// over when it differs from the user-facing [LimitPrice]
	// (e.g. trigger orders adjust the signed price after the
	// trigger fires).
	SignedLimitPrice Decimal `json:"signed_limit_price,omitempty"`
	// SignatureExpiry is the Unix timestamp (seconds) after which the
	// signature is no longer valid.
	SignatureExpiry int64 `json:"signature_expiry_sec,omitempty"`
	// Signature is the EIP-712 action signature attached to the order.
	Signature string `json:"signature,omitempty"`
	// CreationTimestamp is when the engine first saw the order.
	CreationTimestamp MillisTime `json:"creation_timestamp"`
	// LastUpdateTimestamp is the engine's most recent update time.
	LastUpdateTimestamp MillisTime `json:"last_update_timestamp"`

	// Algorithmic-order fields. These are populated only on
	// algo-typed orders (currently TWAP). All four are nullable
	// on the wire.

	// AlgoType identifies the algo strategy ("twap"). Empty for
	// non-algo orders.
	AlgoType string `json:"algo_type,omitempty"`
	// AlgoDurationSec is the algo's total duration in seconds.
	AlgoDurationSec int64 `json:"algo_duration_sec,omitempty"`
	// AlgoNumSlices is the total number of child orders the algo
	// will issue.
	AlgoNumSlices int64 `json:"algo_num_slices,omitempty"`
	// AlgoSlicesCompleted is the number of child orders the algo
	// has placed so far.
	AlgoSlicesCompleted int64 `json:"algo_slices_completed,omitempty"`

	// Trigger-order fields. Populated only on stop-loss /
	// take-profit orders. All five are nullable on the wire.

	// TriggerType identifies the trigger flavour ("stoploss" or
	// "takeprofit"). Empty for non-trigger orders.
	TriggerType string `json:"trigger_type,omitempty"`
	// TriggerPriceType is the price the trigger watches: "mark"
	// or "index".
	TriggerPriceType string `json:"trigger_price_type,omitempty"`
	// TriggerPrice is the price level that activates the order.
	TriggerPrice Decimal `json:"trigger_price,omitempty"`
	// TriggerRejectMessage is the engine's reason if the trigger
	// fired but the resulting order was rejected.
	TriggerRejectMessage string `json:"trigger_reject_message,omitempty"`
}
