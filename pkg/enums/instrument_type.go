// Package enums declares the named-string enums used across the SDK.
//
// Each enum is a defined string type — the simplest idiom in Go that gives
// you exhaustive switch warnings, free JSON round-trips, and a domain-specific
// receiver set without the heavyweight ceremony of an `iota` block plus
// custom marshalers. Aliases of underlying string types like:
//
//	type Direction string
//	const DirectionBuy Direction = "buy"
//
// match what big Go SDKs (aws-sdk-go-v2, stripe-go) use, and the wire format
// they produce is byte-for-byte what Derive expects.
//
// Every enum exposes a Valid method for cheap input validation. Some, like
// [Direction], expose extra domain helpers ([Direction.Sign],
// [Direction.Opposite], [OrderStatus.Terminal]).
//
// # Validating untrusted input
//
// Always check [Direction.Valid] (or the corresponding Valid method on the
// enum) before passing user-provided strings into the SDK. The Go type
// system can't prevent constructing an out-of-range value via `Direction("x")`,
// so the runtime check is the safety net.
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
