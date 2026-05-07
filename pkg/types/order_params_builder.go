// Package types.
package types

import (
	"errors"
	"fmt"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

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
func NewOrderParams(instrument string, side enums.Direction, kind enums.OrderType, amount, limitPrice Decimal) OrderParams {
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

// WithDirection returns a copy with [enums.Direction] set.
func (p OrderParams) WithDirection(d enums.Direction) OrderParams { p.Direction = d; return p }

// WithOrderType returns a copy with [enums.OrderType] set.
func (p OrderParams) WithOrderType(o enums.OrderType) OrderParams { p.OrderType = o; return p }

// WithTimeInForce returns a copy with [enums.TimeInForce] set.
func (p OrderParams) WithTimeInForce(tif enums.TimeInForce) OrderParams {
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
