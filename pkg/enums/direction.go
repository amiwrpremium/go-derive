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

// Direction is the side of a trade or order. Buy orders consume asks; sell
// orders consume bids.
type Direction string

const (
	// DirectionBuy means the order or trade is on the bid side.
	DirectionBuy Direction = "buy"
	// DirectionSell means the order or trade is on the ask side.
	DirectionSell Direction = "sell"
)

// Valid reports whether the receiver equals one of the defined directions.
// Use it to gate untrusted input before [Direction.Sign] or [Direction.Opposite]
// (both of which assume validity).
func (d Direction) Valid() bool {
	switch d {
	case DirectionBuy, DirectionSell:
		return true
	default:
		return false
	}
}

// Sign returns +1 for [DirectionBuy] and -1 for [DirectionSell]. It is
// useful when computing signed position deltas:
//
//	delta := amount.Mul(decimal.NewFromInt(int64(side.Sign())))
//
// Sign panics on values that haven't passed [Direction.Valid]. Validate
// untrusted input first.
func (d Direction) Sign() int {
	switch d {
	case DirectionBuy:
		return 1
	case DirectionSell:
		return -1
	default:
		panic("enums: Direction.Sign called on invalid value " + string(d))
	}
}

// Opposite returns the reverse of d. Used in cancel-and-reverse logic and
// when computing offsetting orders for hedges.
//
// Note: Opposite returns [DirectionBuy] for any non-[DirectionBuy] input,
// including invalid values. Combine with [Direction.Valid] when input
// trustworthiness matters.
func (d Direction) Opposite() Direction {
	if d == DirectionBuy {
		return DirectionSell
	}
	return DirectionBuy
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (d Direction) Validate() error {
	if d.Valid() {
		return nil
	}
	return invalid("Direction", string(d))
}
