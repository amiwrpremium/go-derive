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

import "fmt"

// validationError is returned by every enum's Validate() method when the
// receiver carries a value not in the canonical set. It implements the
// error interface and is comparable with errors.Is against
// [ErrInvalidEnum].
type validationError struct {
	enum  string
	value string
}

func (e *validationError) Error() string {
	return fmt.Sprintf("enums: invalid %s value %q", e.enum, e.value)
}

// Is satisfies errors.Is so callers can match without unwrapping. Every
// enum-validation failure unwraps to [ErrInvalidEnum].
func (e *validationError) Is(target error) bool { return target == ErrInvalidEnum }

// ErrInvalidEnum is the sentinel returned from every enum's Validate
// method when the receiver isn't one of the defined wire values. Use
// errors.Is to detect it.
var ErrInvalidEnum = &validationError{enum: "<unknown>", value: ""}

// invalid is the package-internal helper each Validate method calls to
// build a concrete error if Valid() returned false.
func invalid(enum, value string) error {
	return &validationError{enum: enum, value: value}
}
