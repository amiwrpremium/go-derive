// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the per-asset descriptor returned by
// `public/get_asset` and `public/get_assets`.
package types

import (
	"encoding/json"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

// Asset describes one tradable asset record on Derive — strictly the
// asset entity (the ERC-1155 token), not the "instrument" view that
// adds market/orderbook data.
//
// Mirrors the response shape per
// docs.derive.xyz/reference/public-get_asset.
type Asset struct {
	// Address is the on-chain Asset.sol contract address.
	Address Address `json:"address"`
	// AssetID is the asset's id on the Asset.sol contract.
	AssetID string `json:"asset_id"`
	// AssetName is the canonical asset name.
	AssetName string `json:"asset_name"`
	// AssetType is the asset class — "erc20", "option", or "perp".
	AssetType enums.AssetType `json:"asset_type"`
	// Currency is the underlying currency (e.g. "ETH").
	Currency string `json:"currency"`
	// IsCollateral reports whether the asset can be used as
	// collateral in margin calculations.
	IsCollateral bool `json:"is_collateral"`
	// IsPosition reports whether the asset is treated as a position
	// in margin calculations.
	IsPosition bool `json:"is_position"`
	// ERC20Details is the ERC-20-specific block (decimals,
	// borrow / supply indices, underlying ERC-20 address). Nullable
	// — kept as raw JSON since the subset used varies; decode
	// against `decimals`/`borrow_index`/`supply_index`/
	// `underlying_erc20_address` if needed.
	ERC20Details json.RawMessage `json:"erc20_details,omitempty"`
	// OptionDetails is the option-specific block (expiry, index,
	// option_type, strike, settlement_price). Nullable.
	OptionDetails json.RawMessage `json:"option_details,omitempty"`
	// PerpDetails is the perp-specific block (aggregate_funding,
	// funding rates, etc.). Nullable.
	PerpDetails json.RawMessage `json:"perp_details,omitempty"`
}
