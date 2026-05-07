// Package enums — see asset_type.go for the overview.
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
