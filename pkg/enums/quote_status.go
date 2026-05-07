// Package enums — see asset_type.go for the overview.
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
