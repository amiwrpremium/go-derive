// Package types — see address.go for the overview.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// OrderParams is the request shape for `private/order`.
//
// Most fields map directly to the JSON-RPC schema. The four signing fields
// (Signer, Signature, Nonce, SignatureExpiry) are populated automatically by
// [github.com/amiwrpremium/go-derive/internal/methods.API.PlaceOrder] using
// the configured signer; callers building this struct manually must populate
// them themselves and produce a matching EIP-712 signature.
type OrderParams struct {
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// Direction is buy or sell.
	Direction enums.Direction `json:"direction"`
	// OrderType is limit or market.
	OrderType enums.OrderType `json:"order_type"`
	// TimeInForce is the order's expiry policy.
	TimeInForce enums.TimeInForce `json:"time_in_force,omitempty"`
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
