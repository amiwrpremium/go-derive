// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the user-facing input DTO for `PlaceAlgoOrder`,
// the convenience wrapper around `private/algo_order`.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// AlgoOrderInput parameterises an algorithmic (e.g. TWAP) order.
// Embeds [PlaceOrderInput] for the standard order fields plus the
// algo-specific parameters; the SDK fills in subaccount id,
// signature, signer, nonce and expiry from the configured signer
// and ambient state, identical to PlaceOrder.
//
// AlgoType is the strategy ("twap"). AlgoDurationSec is the total
// execution window in seconds. AlgoNumSlices is the number of child
// executions the parent gets sliced into.
type AlgoOrderInput struct {
	PlaceOrderInput
	AlgoType        enums.AlgoType
	AlgoDurationSec int64
	AlgoNumSlices   int64
}

// Validate performs schema-level checks on the receiver. Wraps the
// embedded [PlaceOrderInput.Validate] with the algo-specific
// requirements (positive duration and slice count, recognised algo
// type).
func (in AlgoOrderInput) Validate() error {
	if err := in.PlaceOrderInput.Validate(); err != nil {
		return err
	}
	if err := in.AlgoType.Validate(); err != nil {
		return invalidParam("algo_type", err.Error())
	}
	if in.AlgoDurationSec <= 0 {
		return invalidParam("algo_duration_sec", "must be positive")
	}
	if in.AlgoNumSlices <= 0 {
		return invalidParam("algo_num_slices", "must be positive")
	}
	return nil
}
