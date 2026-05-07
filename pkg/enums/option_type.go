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

// OptionType distinguishes calls from puts. It only applies when the
// surrounding instrument's [InstrumentType] is [InstrumentTypeOption].
//
// The wire format is single-letter — Derive emits `"C"` for calls and
// `"P"` for puts.
type OptionType string

const (
	// OptionTypeCall gives the holder the right to buy the underlying at
	// the strike on or before expiry.
	OptionTypeCall OptionType = "C"
	// OptionTypePut gives the holder the right to sell the underlying at
	// the strike on or before expiry.
	OptionTypePut OptionType = "P"
)

// Valid reports whether the receiver is one of the defined option types.
func (o OptionType) Valid() bool {
	switch o {
	case OptionTypeCall, OptionTypePut:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (o OptionType) Validate() error {
	if o.Valid() {
		return nil
	}
	return invalid("OptionType", string(o))
}
