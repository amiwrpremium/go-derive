package enums

// LiquidityRole is the role of a fill from the perspective of one side of
// the trade. Maker fills sit on the book first; taker fills cross the book.
//
// Fees and rebates differ by role on every Derive market.
type LiquidityRole string

const (
	// LiquidityRoleMaker means the side provided liquidity (the resting
	// order).
	LiquidityRoleMaker LiquidityRole = "maker"
	// LiquidityRoleTaker means the side consumed liquidity (the crossing
	// order).
	LiquidityRoleTaker LiquidityRole = "taker"
)

// Valid reports whether the receiver is one of the defined roles.
func (r LiquidityRole) Valid() bool {
	switch r {
	case LiquidityRoleMaker, LiquidityRoleTaker:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (r LiquidityRole) Validate() error {
	if r.Valid() {
		return nil
	}
	return invalid("LiquidityRole", string(r))
}
