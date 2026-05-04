package enums

// MarginType identifies which margin model a subaccount uses.
//
// Derive supports three: Standard Margin (`SM`), Portfolio Margin (`PM`),
// and the second-generation Portfolio Margin model (`PM2`). The wire
// values are uppercase abbreviations.
type MarginType string

const (
	// MarginTypeSM is Standard Margin — per-position margin computed
	// independently of the rest of the book. Most permissive accounts.
	MarginTypeSM MarginType = "SM"
	// MarginTypePM is the original Portfolio Margin model — netted
	// margin across the whole subaccount.
	MarginTypePM MarginType = "PM"
	// MarginTypePM2 is the second-generation Portfolio Margin model.
	MarginTypePM2 MarginType = "PM2"
)

// Valid reports whether the receiver is one of the defined margin types.
func (m MarginType) Valid() bool {
	switch m {
	case MarginTypeSM, MarginTypePM, MarginTypePM2:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (m MarginType) Validate() error {
	if m.Valid() {
		return nil
	}
	return invalid("MarginType", string(m))
}
