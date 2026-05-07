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

// AssetType identifies the kind of asset that backs a [Balance] or
// [Collateral] entry. The set is the same as [InstrumentType] — Derive
// reuses the same three categories — but the wire field is named
// `asset_type` rather than `instrument_type` on those payloads.
type AssetType string

const (
	// AssetTypeERC20 is a spot ERC-20 token (typically used as collateral).
	AssetTypeERC20 AssetType = "erc20"
	// AssetTypeOption is an option contract.
	AssetTypeOption AssetType = "option"
	// AssetTypePerp is a perpetual futures contract.
	AssetTypePerp AssetType = "perp"
)

// Valid reports whether the receiver is one of the defined asset types.
func (a AssetType) Valid() bool {
	switch a {
	case AssetTypeERC20, AssetTypeOption, AssetTypePerp:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (a AssetType) Validate() error {
	if a.Valid() {
		return nil
	}
	return invalid("AssetType", string(a))
}
