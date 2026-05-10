// Package methods is the shared implementation of every JSON-RPC method
// Derive exposes. Both pkg/rest.Client and pkg/ws.Client embed *API so that
// each method is defined exactly once, parameterised by the underlying
// transport.
//
// Public methods are unauthenticated; private methods require Signer to be
// non-nil. Private methods that mutate orders also use the Domain to sign
// the per-action EIP-712 hash.
package methods

import (
	"context"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// GetVaultBalances returns one wallet's vault-token holdings. Public.
//
// Optional `params`: `wallet` (smart-contract wallet address) or
// `smart_contract_owner` (EOA that owns the smart-contract wallet).
// At least one of the two must be supplied.
func (a *API) GetVaultBalances(ctx context.Context, params map[string]any) ([]types.VaultBalance, error) {
	if params == nil {
		params = map[string]any{}
	}
	var resp []types.VaultBalance
	if err := a.call(ctx, "public/get_vault_balances", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetVaultShare returns per-block snapshots of one vault token's
// price-per-share over the requested window. Public.
//
// Required `params`: `vault_name`, `from_timestamp_sec`,
// `to_timestamp_sec`. Optional: `page`, `page_size`. Paginated; the
// second return value carries the totals.
func (a *API) GetVaultShare(ctx context.Context, params map[string]any) ([]types.VaultShare, types.Page, error) {
	var resp struct {
		VaultShares []types.VaultShare `json:"vault_shares"`
		Pagination  types.Page         `json:"pagination"`
	}
	if err := a.call(ctx, "public/get_vault_share", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.VaultShares, resp.Pagination, nil
}

// GetVaultAssets lists every ERC-20 asset tracked by Derive's vault
// orderbook. Public.
func (a *API) GetVaultAssets(ctx context.Context) ([]types.VaultAsset, error) {
	var resp []types.VaultAsset
	if err := a.call(ctx, "public/get_vault_assets", map[string]any{}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetVaultStatistics returns a summary snapshot for every Derive
// vault — price-per-share, total supply, USD TVL, and the
// last-trade subaccount value. Public.
func (a *API) GetVaultStatistics(ctx context.Context) ([]types.VaultStatistics, error) {
	var resp []types.VaultStatistics
	if err := a.call(ctx, "public/get_vault_statistics", map[string]any{}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
