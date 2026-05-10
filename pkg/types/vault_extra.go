// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the vault-detail shapes returned by
// `public/get_vault_assets`, `public/get_vault_pools`, and
// `public/get_vault_rates`.
package types

// VaultAsset is one entry in `public/get_vault_assets`. Each entry
// describes one ERC-20 asset that one Derive vault tracks (the
// asset's id, chain, address, and integrator name).
type VaultAsset struct {
	// AssetID is the DB id of the asset.
	AssetID string `json:"asset_id"`
	// ChainID is the chain the asset lives on.
	ChainID int64 `json:"chain_id"`
	// ERC20Address is the on-chain ERC-20 address.
	ERC20Address Address `json:"erc20_address"`
	// Integrator is the integrator name (e.g. the vault provider).
	Integrator string `json:"integrator,omitempty"`
	// Name is the vault asset's display name.
	Name string `json:"name,omitempty"`
	// RPCURL is the chain's RPC endpoint.
	RPCURL string `json:"rpc_url,omitempty"`
}

// VaultPool is one entry in `public/get_vault_pools` — a registered
// vault ERC-20 pool (the manager contract that holds vault deposits).
type VaultPool struct {
	// Address is the manager-contract address.
	Address Address `json:"address"`
	// ChainID is the chain the pool lives on.
	ChainID int64 `json:"chain_id"`
	// Name is the pool's display name.
	Name string `json:"name"`
	// PoolType is the pool's category (e.g. "basis", "perp").
	PoolType string `json:"pool_type,omitempty"`
}

// VaultRates is the response of `public/get_vault_rates` — the
// engine's current view of one vault's basis / interest / funding
// economics.
type VaultRates struct {
	// Rate is the headline yield rate on the vault.
	Rate Decimal `json:"rate"`
	// TotalRate is the aggregated rate across components.
	TotalRate Decimal `json:"total_rate,omitempty"`
	// FundingRate is the underlying perp funding rate the vault
	// captures.
	FundingRate Decimal `json:"funding_rate,omitempty"`
	// InterestRate is the lending-side rate the vault picks up.
	InterestRate Decimal `json:"interest_rate,omitempty"`
	// LRTPrice is the LRT vault token's reference price.
	LRTPrice Decimal `json:"lrt_price,omitempty"`
	// BaseBalance is the vault's base-asset balance.
	BaseBalance Decimal `json:"base_balance,omitempty"`
	// QuoteBalance is the vault's quote-asset balance.
	QuoteBalance Decimal `json:"quote_balance,omitempty"`
	// PerpBalance is the vault's perp-leg balance.
	PerpBalance Decimal `json:"perp_balance,omitempty"`
	// USDDepositAmount is the vault's $ deposit amount.
	USDDepositAmount Decimal `json:"usd_deposit_amount,omitempty"`
	// YearlyFunding is the projected annual funding-component yield.
	YearlyFunding Decimal `json:"yearly_funding,omitempty"`
	// YearlyInterest is the projected annual interest-component yield.
	YearlyInterest Decimal `json:"yearly_interest,omitempty"`
	// YearlyStaticLRTYield is the projected annual static LRT yield.
	YearlyStaticLRTYield Decimal `json:"yearly_static_LRT_yield,omitempty"`
}
