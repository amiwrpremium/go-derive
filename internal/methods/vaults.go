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
// Pass either `wallet` (the smart-contract wallet address) or
// `smartContractOwner` (the EOA that owns the smart-contract
// wallet); at least one must be non-empty.
func (a *API) GetVaultBalances(ctx context.Context, q types.VaultBalancesQuery) ([]types.VaultBalance, error) {
	params := map[string]any{}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	}
	if q.SmartContractOwner != "" {
		params["smart_contract_owner"] = q.SmartContractOwner
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
// Timestamps on this endpoint are in seconds, not milliseconds —
// match the field names on [types.VaultShareQuery]. Paginated; the
// second return value carries the totals.
func (a *API) GetVaultShare(ctx context.Context, q types.VaultShareQuery, page types.PageRequest) ([]types.VaultShare, types.Page, error) {
	params := map[string]any{
		"vault_name":         q.VaultName,
		"from_timestamp_sec": q.FromSec,
		"to_timestamp_sec":   q.ToSec,
	}
	addPaging(params, page)
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

// GetVaultPools lists every registered vault ERC-20 pool — the
// manager contracts that hold vault deposits. Public.
func (a *API) GetVaultPools(ctx context.Context) ([]types.VaultPool, error) {
	var resp []types.VaultPool
	if err := a.call(ctx, "public/get_vault_pools", map[string]any{}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetVaultRates returns the engine's current view of one basis
// vault's rate components — funding, interest, LRT yield, and
// the per-leg balances. Public.
//
// Required `vaultType`: documented values are `lbtc` and `weeth`.
func (a *API) GetVaultRates(ctx context.Context, q types.VaultRatesQuery) (types.VaultRates, error) {
	params := map[string]any{}
	if q.VaultType != "" {
		params["vault_type"] = q.VaultType
	}
	var resp types.VaultRates
	if err := a.call(ctx, "public/get_vault_rates", params, &resp); err != nil {
		return types.VaultRates{}, err
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
