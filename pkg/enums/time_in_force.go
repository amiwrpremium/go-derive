// Package enums — see asset_type.go for the overview.
package enums

// TimeInForce controls when, and under what conditions, an order becomes
// inactive. The matching engine consults the time-in-force as soon as the
// order is accepted; a TIF mismatch (e.g. PostOnly on a market order) yields
// a synchronous rejection.
type TimeInForce string

const (
	// TimeInForceGTC ("good-till-cancelled") keeps the order open until the
	// caller cancels it or it expires for another reason.
	TimeInForceGTC TimeInForce = "gtc"
	// TimeInForcePostOnly rejects the order if it would cross the book at
	// submission time. Used by makers to guarantee maker rebates.
	TimeInForcePostOnly TimeInForce = "post_only"
	// TimeInForceFOK ("fill-or-kill") requires the order to fill in full
	// immediately or be cancelled. Partial fills are not allowed.
	TimeInForceFOK TimeInForce = "fok"
	// TimeInForceIOC ("immediate-or-cancel") fills as much as it can right
	// now and cancels any remaining quantity. Partial fills are allowed.
	TimeInForceIOC TimeInForce = "ioc"
)

// Valid reports whether the receiver is one of the defined TIFs.
func (t TimeInForce) Valid() bool {
	switch t {
	case TimeInForceGTC, TimeInForcePostOnly, TimeInForceFOK, TimeInForceIOC:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (t TimeInForce) Validate() error {
	if t.Valid() {
		return nil
	}
	return invalid("TimeInForce", string(t))
}
