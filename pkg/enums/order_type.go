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

// OrderType describes how an order interacts with the order book.
//
// [OrderTypeLimit] orders rest on the book at a stated price; [OrderTypeMarket]
// orders cross the book immediately at the best available price subject to
// the user's slippage cap.
type OrderType string

const (
	// OrderTypeLimit is a price-limited order that rests on the book until
	// it crosses, expires, or is cancelled.
	OrderTypeLimit OrderType = "limit"
	// OrderTypeMarket is an order that crosses the book immediately,
	// constrained only by the caller's max-fee and limit-price cap.
	OrderTypeMarket OrderType = "market"
)

// Valid reports whether the receiver is one of the defined order types.
func (t OrderType) Valid() bool {
	switch t {
	case OrderTypeLimit, OrderTypeMarket:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (t OrderType) Validate() error {
	if t.Valid() {
		return nil
	}
	return invalid("OrderType", string(t))
}
