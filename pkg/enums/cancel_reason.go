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
		CancelReasonSubaccountWithdrawn, CancelReasonCompliance:
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
