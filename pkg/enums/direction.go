// Package enums — see asset_type.go for the overview.
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
