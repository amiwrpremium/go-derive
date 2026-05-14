// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the user-facing input DTO for `PlaceOrder` —
// previously declared in internal/methods, lifted here so callers
// only need to import pkg/types for the SDK's domain types.
package types

import (
	"github.com/amiwrpremium/go-derive/pkg/enums"
)

// PlaceOrderInput is a thin convenience wrapper for the user-facing
// PlaceOrder method. It contains only the strategically-relevant fields;
// the SDK fills in subaccount id, signature, signer, nonce and expiry from
// the configured signer and ambient state.
type PlaceOrderInput struct {
	InstrumentName string
	Asset          Address
	SubID          uint64
	Direction      enums.Direction
	OrderType      enums.OrderType
	TimeInForce    enums.TimeInForce
	Amount         Decimal
	LimitPrice     Decimal
	MaxFee         Decimal
	Label          string
	MMP            bool
	ReduceOnly     bool
	// Client is an optional caller-defined tag attached to the order.
	// Echoed back on the order record; useful for client-side
	// reconciliation.
	Client string
	// IsAtomicSigning marks the order as signed via EIP-1271
	// (Safe / smart-contract wallet). Required by callers using a
	// vault-style signer rather than a plain EOA.
	IsAtomicSigning bool
	// ReferralCode is an optional referral code applied to the order.
	ReferralCode string
	// RejectPostOnly rejects the order if it would have crossed the
	// book (stricter form of post-only — fails fast instead of
	// resting).
	RejectPostOnly bool
	// RejectTimestamp is the latest acceptable arrival time at the
	// engine, in milliseconds since the Unix epoch. The engine
	// rejects the order if its own clock is past this value.
	RejectTimestamp int64
	// ExtraFee is an optional caller-paid tip on top of the standard
	// fee schedule. Denominated in quote currency (USDC).
	ExtraFee Decimal
}

// Validate performs schema-level checks on the receiver: required fields
// populated, enum values in range, numeric fields in bounds. It does not
// validate against an instrument's tick / amount step (those live on
// [Instrument] and require a network round-trip).
//
// Returns nil on success or an error wrapping [ErrInvalidParams].
func (in PlaceOrderInput) Validate() error {
	if in.InstrumentName == "" {
		return invalidParam("instrument_name", "required")
	}
	if in.Asset.IsZero() {
		return invalidParam("asset", "required")
	}
	if err := in.Direction.Validate(); err != nil {
		return invalidParam("direction", err.Error())
	}
	if err := in.OrderType.Validate(); err != nil {
		return invalidParam("order_type", err.Error())
	}
	if in.TimeInForce != "" {
		if err := in.TimeInForce.Validate(); err != nil {
			return invalidParam("time_in_force", err.Error())
		}
	}
	if in.Amount.Sign() <= 0 {
		return invalidParam("amount", "must be positive")
	}
	if in.LimitPrice.Sign() <= 0 {
		return invalidParam("limit_price", "must be positive")
	}
	if in.MaxFee.Sign() < 0 {
		return invalidParam("max_fee", "must be non-negative")
	}
	return nil
}
