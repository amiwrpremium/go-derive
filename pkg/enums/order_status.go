package enums

// OrderStatus is the lifecycle state of an order as reported by the
// matching engine. Use [OrderStatus.Terminal] to test for "no further
// updates expected".
//
// The set mirrors the canonical `derivexyz/cockpit` enum exactly:
// open, filled, cancelled, expired, rejected.
type OrderStatus string

const (
	// OrderStatusOpen means the order is resting on the book.
	OrderStatusOpen OrderStatus = "open"
	// OrderStatusFilled means the order has been completely matched.
	OrderStatusFilled OrderStatus = "filled"
	// OrderStatusCancelled means the order was cancelled by the user, the
	// session-key, or the engine before it filled. The associated
	// [CancelReason] explains which.
	OrderStatusCancelled OrderStatus = "cancelled"
	// OrderStatusExpired means the order's signature expiry passed before
	// it filled.
	OrderStatusExpired OrderStatus = "expired"
	// OrderStatusRejected means the engine rejected the order at submission
	// time (e.g. invalid price, post-only would cross).
	OrderStatusRejected OrderStatus = "rejected"
)

// Valid reports whether the receiver is one of the defined statuses.
func (s OrderStatus) Valid() bool {
	switch s {
	case OrderStatusOpen, OrderStatusFilled, OrderStatusCancelled,
		OrderStatusExpired, OrderStatusRejected:
		return true
	default:
		return false
	}
}

// Terminal reports whether the status is final — i.e. the order will not
// receive further updates and can be cleaned out of any in-memory cache.
//
// Only Open is non-terminal; everything else is.
func (s OrderStatus) Terminal() bool {
	switch s {
	case OrderStatusFilled, OrderStatusCancelled, OrderStatusExpired,
		OrderStatusRejected:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (s OrderStatus) Validate() error {
	if s.Valid() {
		return nil
	}
	return invalid("OrderStatus", string(s))
}
