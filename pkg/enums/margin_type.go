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

// MarginType identifies which margin model a subaccount uses.
//
// Derive supports three: Standard Margin (`SM`), Portfolio Margin (`PM`),
// and the second-generation Portfolio Margin model (`PM2`). The wire
// values are uppercase abbreviations.
type MarginType string

const (
	// MarginTypeSM is Standard Margin — per-position margin computed
	// independently of the rest of the book. Most permissive accounts.
	MarginTypeSM MarginType = "SM"
	// MarginTypePM is the original Portfolio Margin model — netted
	// margin across the whole subaccount.
	MarginTypePM MarginType = "PM"
	// MarginTypePM2 is the second-generation Portfolio Margin model.
	MarginTypePM2 MarginType = "PM2"
)

// Valid reports whether the receiver is one of the defined margin types.
func (m MarginType) Valid() bool {
	switch m {
	case MarginTypeSM, MarginTypePM, MarginTypePM2:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (m MarginType) Validate() error {
	if m.Valid() {
		return nil
	}
	return invalid("MarginType", string(m))
}
