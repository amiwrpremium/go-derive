// Package enums declares the named-string enums used across the SDK.
//
// This file holds [AlgoType], the algorithmic-order strategy submitted
// to `private/algo_order`.
package enums

// AlgoType is the algorithmic-order strategy. Today only TWAP is
// documented; the SDK uses a defined-string type so future strategies
// (VWAP, etc.) extend the enum without changing field types.
type AlgoType string

const (
	// AlgoTypeTWAP slices the parent order into equal child executions
	// across [AlgoOrderInput.AlgoDurationSec] seconds.
	AlgoTypeTWAP AlgoType = "twap"
)

// Valid reports whether the receiver is one of the defined algo types.
func (a AlgoType) Valid() bool {
	switch a {
	case AlgoTypeTWAP:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (a AlgoType) Validate() error {
	if a.Valid() {
		return nil
	}
	return invalid("AlgoType", string(a))
}
