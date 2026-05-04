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
