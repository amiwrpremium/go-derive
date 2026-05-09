// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shapes for the public vault endpoints:
// `public/get_vault_balances`, `public/get_vault_share`, and
// `public/get_vault_statistics`.
package types

// VaultBalance is one entry in `public/get_vault_balances`. Each entry
// reports a wallet's holding of one Derive vault token.
//
// The shape mirrors `VaultBalanceResponseSchema` in Derive's v2.2
// OpenAPI spec.
type VaultBalance struct {
	// Name is the vault's name.
	Name string `json:"name"`
	// Address is the vault token's contract address.
	Address Address `json:"address"`
	// ChainID is the chain the holding lives on.
	ChainID int64 `json:"chain_id"`
	// VaultAssetType identifies the vault's asset class (e.g.
	// "rswETHC").
	VaultAssetType string `json:"vault_asset_type"`
	// Amount is the wallet's vault-token balance (in the vault
	// token's own units).
	Amount Decimal `json:"amount"`
}

// VaultShare is one entry in `public/get_vault_share.vault_shares`.
// Each entry is a per-block snapshot of the vault token's
// price-per-share against base / underlying / USD.
//
// The shape mirrors `VaultShareResponseSchema`.
type VaultShare struct {
	// BlockNumber is the Derive-chain block this snapshot was taken
	// at.
	BlockNumber int64 `json:"block_number"`
	// BlockTimestamp is the block's timestamp (Unix seconds).
	BlockTimestamp int64 `json:"block_timestamp"`
	// BaseValue is the vault token's price against the vault's base
	// currency (e.g. rswETHC vs rswETH).
	BaseValue Decimal `json:"base_value"`
	// UnderlyingValue is the price against the underlying currency
	// (e.g. rswETHC vs ETH). Nullable on the wire — zero-value
	// Decimal when the vault has no defined underlying.
	UnderlyingValue Decimal `json:"underlying_value,omitempty"`
	// USDValue is the price in USD.
	USDValue Decimal `json:"usd_value"`
}

// VaultStatistics is one entry in `public/get_vault_statistics`. It
// reports the at-block snapshot of one vault's price-per-share, total
// supply, TVL, and the last-trade subaccount value.
//
// The shape mirrors `VaultStatisticsResponseSchema`.
type VaultStatistics struct {
	// VaultName is the vault's name.
	VaultName string `json:"vault_name"`
	// BlockNumber is the Derive-chain block this snapshot was taken
	// at.
	BlockNumber int64 `json:"block_number"`
	// BlockTimestamp is the block's timestamp (Unix seconds).
	BlockTimestamp int64 `json:"block_timestamp"`
	// TotalSupply is the total supply of the vault's token on the
	// Derive chain.
	TotalSupply Decimal `json:"total_supply"`
	// USDTVL is the total USD TVL of the vault.
	USDTVL Decimal `json:"usd_tvl"`
	// USDValue is the vault token's price in USD.
	USDValue Decimal `json:"usd_value"`
	// BaseValue is the vault token's price against the base currency.
	BaseValue Decimal `json:"base_value"`
	// UnderlyingValue is the price against the underlying currency.
	// Nullable on the wire — zero when the vault has no underlying.
	UnderlyingValue Decimal `json:"underlying_value,omitempty"`
	// SubaccountValueAtLastTrade is the vault subaccount's equity at
	// its most recent trade. Nullable on the wire — zero before the
	// vault has subaccount activity.
	SubaccountValueAtLastTrade Decimal `json:"subaccount_value_at_last_trade,omitempty"`
}
