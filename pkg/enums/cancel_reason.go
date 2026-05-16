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

// CancelReason carries the engine's explanation for why an order was
// cancelled. It is reported on the canonical Order record (and on quote
// updates) — `""` for orders that have not been cancelled, otherwise one
// of the values below.
//
// The set mirrors the canonical `derivexyz/cockpit` enum exactly.
type CancelReason string

const (
	// CancelReasonNone is the empty wire value used while the order is
	// still open (no cancel has happened).
	CancelReasonNone CancelReason = ""
	// CancelReasonUserRequest means the order was cancelled by an explicit
	// `private/cancel` (or label/instrument/all variant).
	CancelReasonUserRequest CancelReason = "user_request"
	// CancelReasonMMP means market-maker protection tripped and pulled
	// the order along with the rest of the maker's book.
	CancelReasonMMP CancelReason = "mmp_trigger"
	// CancelReasonInsufficientMargin means the engine cancelled the order
	// because filling it would breach the subaccount's margin rules.
	CancelReasonInsufficientMargin CancelReason = "insufficient_margin"
	// CancelReasonSignedMaxFeeTooLow means the signed `max_fee` field was
	// below the venue's required minimum at the time of fill.
	CancelReasonSignedMaxFeeTooLow CancelReason = "signed_max_fee_too_low"
	// CancelReasonIOC means an IOC or market order partially filled and
	// the remainder was cancelled by IOC semantics.
	CancelReasonIOC CancelReason = "ioc_or_market_partial_fill"
	// CancelReasonCancelOnDisconnect means the kill-switch fired because
	// the wallet's authenticated WebSocket session disconnected.
	CancelReasonCancelOnDisconnect CancelReason = "cancel_on_disconnect"
	// CancelReasonSessionKey means the signing session key was deregistered,
	// which invalidates all of its outstanding orders.
	CancelReasonSessionKey CancelReason = "session_key_deregistered"
	// CancelReasonSubaccountWithdrawn means the subaccount was withdrawn
	// from the venue, taking its outstanding orders with it.
	CancelReasonSubaccountWithdrawn CancelReason = "subaccount_withdrawn"
	// CancelReasonCompliance means the wallet was placed in a restricted
	// compliance state and its open orders were pulled.
	CancelReasonCompliance CancelReason = "compliance"
	// CancelReasonTriggerFailed means a trigger order activated but the
	// child order the trigger placed could not be created (e.g. the
	// engine rejected the resulting fill price).
	CancelReasonTriggerFailed CancelReason = "trigger_failed"
	// CancelReasonValidationFailed means the engine rejected the order
	// during a validation step that doesn't map to one of the more
	// specific reasons.
	CancelReasonValidationFailed CancelReason = "validation_failed"
	// CancelReasonAlgoCompleted means an algo order (currently TWAP)
	// finished its planned slices and its parent order was closed.
	CancelReasonAlgoCompleted CancelReason = "algo_completed"
)

// Valid reports whether the receiver is one of the defined cancel reasons.
//
// `CancelReasonNone` (the empty string) counts as valid — that is the
// wire value for "still open, never cancelled".
func (c CancelReason) Valid() bool {
	switch c {
	case CancelReasonNone, CancelReasonUserRequest, CancelReasonMMP,
		CancelReasonInsufficientMargin, CancelReasonSignedMaxFeeTooLow,
		CancelReasonIOC, CancelReasonCancelOnDisconnect, CancelReasonSessionKey,
		CancelReasonSubaccountWithdrawn, CancelReasonCompliance,
		CancelReasonTriggerFailed, CancelReasonValidationFailed,
		CancelReasonAlgoCompleted:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (c CancelReason) Validate() error {
	if c.Valid() {
		return nil
	}
	return invalid("CancelReason", string(c))
}
