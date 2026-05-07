// Package enums.
package enums

// AuctionType describes the regime of a liquidation auction Derive ran
// against an undercollateralised subaccount.
type AuctionType string

const (
	// AuctionTypeSolvent — subaccount equity is positive; auction
	// transfers positions to keep things solvent.
	AuctionTypeSolvent AuctionType = "solvent"
	// AuctionTypeInsolvent — subaccount equity is negative; insurance
	// fund / socialised-loss path activates.
	AuctionTypeInsolvent AuctionType = "insolvent"
)

// Valid reports whether the receiver is one of the defined auction types.
func (a AuctionType) Valid() bool {
	switch a {
	case AuctionTypeSolvent, AuctionTypeInsolvent:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (a AuctionType) Validate() error {
	if a.Valid() {
		return nil
	}
	return invalid("AuctionType", string(a))
}
