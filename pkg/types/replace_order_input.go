// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the user-facing input DTO for the atomic
// cancel-and-place `Replace` method, wrapping `private/replace`.
package types

// ReplaceOrderInput parameterises an atomic cancel-and-place. The
// embedded [PlaceOrderInput] carries the replacement order's full
// shape; OrderIDToCancel (or NonceToCancel) identifies the order
// being cancelled. The SDK fills in subaccount id, signature,
// signer, nonce, and expiry for the replacement, identical to
// [PlaceOrderInput].
//
// Exactly one of OrderIDToCancel and NonceToCancel must be set.
type ReplaceOrderInput struct {
	PlaceOrderInput
	// OrderIDToCancel is the engine-assigned id of the order to
	// cancel. Set this when you have the id back from a prior
	// PlaceOrder.
	OrderIDToCancel string
	// NonceToCancel is the signed nonce of the order to cancel.
	// Set this when cancelling an order whose id you have not
	// received yet (e.g. cancel-on-disconnect race).
	NonceToCancel uint64
}

// Validate performs schema-level checks on the receiver. Wraps the
// embedded [PlaceOrderInput.Validate] with the cancel-target
// requirement (exactly one of OrderIDToCancel / NonceToCancel set).
func (in ReplaceOrderInput) Validate() error {
	if err := in.PlaceOrderInput.Validate(); err != nil {
		return err
	}
	hasID := in.OrderIDToCancel != ""
	hasNonce := in.NonceToCancel != 0
	if !hasID && !hasNonce {
		return invalidParam("order_id_to_cancel|nonce_to_cancel", "one of order_id_to_cancel or nonce_to_cancel is required")
	}
	if hasID && hasNonce {
		return invalidParam("order_id_to_cancel|nonce_to_cancel", "must not set both order_id_to_cancel and nonce_to_cancel")
	}
	return nil
}
