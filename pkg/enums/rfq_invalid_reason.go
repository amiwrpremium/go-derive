// Package enums declares the named-string enums used across the SDK.
//
// Each enum is a defined string type — the simplest idiom in Go that gives
// you exhaustive switch warnings, free JSON round-trips, and a domain-specific
// receiver set without the heavyweight ceremony of an `iota` block plus
// custom marshalers. Aliases of underlying string types like:
//
//	type Direction string
//	const DirectionBuy Direction = "buy"
//
// match what big Go SDKs (aws-sdk-go-v2, stripe-go) use, and the wire format
// they produce is byte-for-byte what Derive expects.
//
// Every enum exposes a Valid method for cheap input validation. Some, like
// [Direction], expose extra domain helpers ([Direction.Sign],
// [Direction.Opposite], [OrderStatus.Terminal]).
//
// # Validating untrusted input
//
// Always check [Direction.Valid] (or the corresponding Valid method on the
// enum) before passing user-provided strings into the SDK. The Go type
// system can't prevent constructing an out-of-range value via `Direction("x")`,
// so the runtime check is the safety net.
package enums

// RFQInvalidReason is the engine's reason for rejecting an RFQ /
// order-quote pre-flight as invalid (margin or buying-power
// constraints).
//
// The wire field is nullable — a nil/empty value means the request is
// valid and there is no reason to surface. The seven non-empty values
// below mirror the `invalid_reason` enum on
// `PrivateRfqGetBestQuoteResultSchema` verbatim, including punctuation
// and casing, so they round-trip as-is.
type RFQInvalidReason string

const (
	// RFQInvalidReasonAccountUnderMaintenance — the account is already
	// under its maintenance-margin floor; all trading is frozen.
	RFQInvalidReasonAccountUnderMaintenance RFQInvalidReason = "Account is currently under maintenance margin requirements, trading is frozen."
	// RFQInvalidReasonWouldUnderMaintenance — placing the trade would
	// drop the account below its maintenance-margin floor.
	RFQInvalidReasonWouldUnderMaintenance RFQInvalidReason = "This order would cause account to fall under maintenance margin requirements."
	// RFQInvalidReasonRiskReducingOnly — the account is restricted to a
	// single risk-reducing open order at a time.
	RFQInvalidReasonRiskReducingOnly RFQInvalidReason = "Insufficient buying power, only a single risk-reducing open order is allowed."
	// RFQInvalidReasonReduceSize — the trade exceeds buying power; the
	// engine suggests reducing order size.
	RFQInvalidReasonReduceSize RFQInvalidReason = "Insufficient buying power, consider reducing order size."
	// RFQInvalidReasonReduceOrCancel — the trade exceeds buying power;
	// the engine suggests reducing order size or canceling other
	// orders.
	RFQInvalidReasonReduceOrCancel RFQInvalidReason = "Insufficient buying power, consider reducing order size or canceling other orders."
	// RFQInvalidReasonCancelLimitsOrUseIOC — the trade is risk-reducing
	// in isolation but interacts with other open orders such that
	// buying power might be insufficient on a partial fill.
	RFQInvalidReasonCancelLimitsOrUseIOC RFQInvalidReason = "Consider canceling other limit orders or using IOC, FOK, or market orders. This order is risk-reducing, but if filled with other open orders, buying power might be insufficient."
	// RFQInvalidReasonInsufficientBuyingPower — generic insufficient
	// buying power; no actionable suggestion attached.
	RFQInvalidReasonInsufficientBuyingPower RFQInvalidReason = "Insufficient buying power."
)

// Valid reports whether the receiver is one of the seven documented
// invalid-reason strings. The empty string (the "no reason" wire
// value) is not considered Valid by this method — if you want to
// distinguish "valid request, no reason" from "unknown reason", check
// for the empty string explicitly.
func (r RFQInvalidReason) Valid() bool {
	switch r {
	case RFQInvalidReasonAccountUnderMaintenance,
		RFQInvalidReasonWouldUnderMaintenance,
		RFQInvalidReasonRiskReducingOnly,
		RFQInvalidReasonReduceSize,
		RFQInvalidReasonReduceOrCancel,
		RFQInvalidReasonCancelLimitsOrUseIOC,
		RFQInvalidReasonInsufficientBuyingPower:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the documented
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (r RFQInvalidReason) Validate() error {
	if r.Valid() {
		return nil
	}
	return invalid("RFQInvalidReason", string(r))
}
