package enums

// InstrumentType identifies the kind of contract a market quotes.
//
// Derive supports three: linear perpetuals, options (with expiries and
// strikes), and ERC-20 spot tokens used as collateral or for spot trading.
type InstrumentType string

const (
	// InstrumentTypePerp is a perpetual futures contract with continuous
	// funding payments and no fixed expiry.
	InstrumentTypePerp InstrumentType = "perp"
	// InstrumentTypeOption is a European-style option with a strike and
	// expiry; see [github.com/amiwrpremium/go-derive/pkg/types.OptionDetails].
	InstrumentTypeOption InstrumentType = "option"
	// InstrumentTypeERC20 is a spot ERC-20 token (typically used as collateral).
	InstrumentTypeERC20 InstrumentType = "erc20"
)

// Valid reports whether the receiver is one of the defined instrument types.
func (k InstrumentType) Valid() bool {
	switch k {
	case InstrumentTypePerp, InstrumentTypeOption, InstrumentTypeERC20:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (k InstrumentType) Validate() error {
	if k.Valid() {
		return nil
	}
	return invalid("InstrumentType", string(k))
}
