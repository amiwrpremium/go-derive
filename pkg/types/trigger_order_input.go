// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the user-facing input DTO for `PlaceTriggerOrder`,
// the convenience wrapper around `private/trigger_order`.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// TriggerOrderInput parameterises a stop-loss or take-profit order.
// Embeds [PlaceOrderInput] for the standard order fields plus the
// trigger-specific parameters; the SDK fills in subaccount id,
// signature, signer, nonce and expiry from the configured signer
// and ambient state, identical to PlaceOrder.
//
// TriggerType is the flavour (stoploss or takeprofit).
// TriggerPriceType is which price the engine watches (mark or
// index). TriggerPrice is the level at which the order activates;
// once active the resulting order behaves as a regular limit /
// market order at LimitPrice.
type TriggerOrderInput struct {
	PlaceOrderInput
	TriggerType      enums.TriggerType
	TriggerPriceType enums.TriggerPriceType
	TriggerPrice     Decimal
}

// Validate performs schema-level checks on the receiver. Wraps the
// embedded [PlaceOrderInput.Validate] with the trigger-specific
// requirements (recognised trigger and price-type values, positive
// trigger price).
func (in TriggerOrderInput) Validate() error {
	if err := in.PlaceOrderInput.Validate(); err != nil {
		return err
	}
	if err := in.TriggerType.Validate(); err != nil {
		return invalidParam("trigger_type", err.Error())
	}
	if err := in.TriggerPriceType.Validate(); err != nil {
		return invalidParam("trigger_price_type", err.Error())
	}
	if in.TriggerPrice.Sign() <= 0 {
		return invalidParam("trigger_price", "must be positive")
	}
	return nil
}
