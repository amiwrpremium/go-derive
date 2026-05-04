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
