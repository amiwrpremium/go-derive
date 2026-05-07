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

// QuoteStatus is the lifecycle state of a maker [Quote] response to an
// open RFQ.
//
// The set is the same four values [OrderStatus] supports — open / filled
// / cancelled / expired — but exists as its own type so a `Quote.Status`
// field cannot be confused with an `Order.OrderStatus` at the type
// level.
type QuoteStatus string

const (
	// QuoteStatusOpen — quote is live and may still be executed.
	QuoteStatusOpen QuoteStatus = "open"
	// QuoteStatusFilled — taker executed the quote.
	QuoteStatusFilled QuoteStatus = "filled"
	// QuoteStatusCancelled — quote was cancelled (by maker or engine).
	QuoteStatusCancelled QuoteStatus = "cancelled"
	// QuoteStatusExpired — quote's signature_expiry_sec passed.
	QuoteStatusExpired QuoteStatus = "expired"
)

// Valid reports whether the receiver is one of the defined statuses.
func (s QuoteStatus) Valid() bool {
	switch s {
	case QuoteStatusOpen, QuoteStatusFilled, QuoteStatusCancelled, QuoteStatusExpired:
		return true
	default:
		return false
	}
}

// Terminal reports whether the status is final.
func (s QuoteStatus) Terminal() bool {
	switch s {
	case QuoteStatusFilled, QuoteStatusCancelled, QuoteStatusExpired:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (s QuoteStatus) Validate() error {
	if s.Valid() {
		return nil
	}
	return invalid("QuoteStatus", string(s))
}
