// Package types.
package types

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

// Decimal is a fixed-precision decimal number that JSON-encodes as a string.
//
// Derive returns numbers (prices, sizes, fees) as JSON strings to avoid the
// truncation float64 would impose on 18-decimal-place values. [Decimal]
// preserves the round-trip byte-for-byte and supports the full
// shopspring/decimal arithmetic API via [Decimal.Inner].
//
// The zero value is the decimal zero; it is safe to use without
// initialisation.
type Decimal struct{ d decimal.Decimal }

// NewDecimal parses the canonical decimal representation s into a [Decimal].
// Acceptable forms include "0", "1.5", "0.0000000000000000018", "-2.5", and
// scientific notation ("1.5e3"). It returns an error if s is not a valid
// decimal.
func NewDecimal(s string) (Decimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Decimal{}, fmt.Errorf("types: parse decimal %q: %w", s, err)
	}
	return Decimal{d: d}, nil
}

// MustDecimal is [NewDecimal] that panics on failure. It is appropriate in
// constants and tests where the input is known-good and a parse failure
// is a programmer bug.
func MustDecimal(s string) Decimal {
	d, err := NewDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

// DecimalFromInt builds a [Decimal] from a signed integer. It is exact for
// any int64 value.
func DecimalFromInt(n int64) Decimal {
	return Decimal{d: decimal.NewFromInt(n)}
}

// String returns the canonical decimal representation — the same string
// shopspring/decimal would produce, with trailing zeroes stripped.
func (d Decimal) String() string { return d.d.String() }

// Inner returns the underlying shopspring/decimal.Decimal, allowing callers
// to perform arithmetic without a round-trip through string.
//
// The returned value is a copy and is independent of the receiver.
func (d Decimal) Inner() decimal.Decimal { return d.d }

// IsZero reports whether the decimal equals zero.
func (d Decimal) IsZero() bool { return d.d.IsZero() }

// Sign returns -1, 0 or +1 for negative, zero or positive values
// respectively.
func (d Decimal) Sign() int { return d.d.Sign() }

// MarshalJSON encodes the decimal as a JSON string — the form Derive
// expects on the wire.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.d.String())
}

// UnmarshalJSON decodes a JSON string or number into the receiver.
//
// The Derive API always emits strings, but the implementation tolerates
// numeric input for resilience. Empty strings and JSON null leave the
// receiver untouched.
func (d *Decimal) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			return nil
		}
		v, err := decimal.NewFromString(s)
		if err != nil {
			return fmt.Errorf("types: decode decimal %q: %w", s, err)
		}
		d.d = v
		return nil
	}
	v, err := decimal.NewFromString(string(b))
	if err != nil {
		return fmt.Errorf("types: decode decimal %s: %w", b, err)
	}
	d.d = v
	return nil
}
