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
