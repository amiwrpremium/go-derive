package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestGetVaultBalances_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_vault_balances", []any{
		map[string]any{
			"name": "rswETHC", "address": "0x1111111111111111111111111111111111111111",
			"chain_id": int64(1), "vault_asset_type": "rswETHC", "amount": "12.5",
		},
	})
	got, err := api.GetVaultBalances(context.Background(), "0x2222222222222222222222222222222222222222", "")
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "rswETHC", got[0].Name)
	assert.Equal(t, "12.5", got[0].Amount.String())
}

func TestGetVaultShare_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_vault_share", map[string]any{
		"vault_shares": []any{
			map[string]any{
				"block_number": int64(12345), "block_timestamp": int64(1700000000),
				"base_value": "1.05", "underlying_value": "1.07", "usd_value": "2700",
			},
		},
		"pagination": map[string]any{"num_pages": 1, "count": 1},
	})
	shares, page, err := api.GetVaultShare(context.Background(), types.VaultShareQuery{
		VaultName: "rswETHC",
		FromSec:   1700000000,
		ToSec:     1700100000,
	}, types.PageRequest{})
	require.NoError(t, err)
	require.Len(t, shares, 1)
	assert.Equal(t, "1.07", shares[0].UnderlyingValue.String())
	assert.Equal(t, 1, page.Count)
}

func TestGetVaultAssets_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_vault_assets", []any{
		map[string]any{
			"asset_id": "1", "chain_id": int64(1),
			"erc20_address": "0x1111111111111111111111111111111111111111",
			"integrator":    "across", "name": "rswETH", "rpc_url": "https://rpc.example",
		},
	})
	got, err := api.GetVaultAssets(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "rswETH", got[0].Name)
}

func TestGetVaultPools_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_vault_pools", []any{
		map[string]any{
			"address":  "0x1111111111111111111111111111111111111111",
			"chain_id": int64(1), "name": "rswETH-pool", "pool_type": "basis",
		},
	})
	got, err := api.GetVaultPools(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "rswETH-pool", got[0].Name)
}

func TestGetVaultRates_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_vault_rates", map[string]any{
		"rate":           "0.05",
		"total_rate":     "0.07",
		"funding_rate":   "0.01",
		"interest_rate":  "0.04",
		"lrt_price":      "1.05",
		"base_balance":   "100",
		"quote_balance":  "200",
		"perp_balance":   "0",
		"yearly_funding": "0.10",
	})
	got, err := api.GetVaultRates(context.Background(), "weeth")
	require.NoError(t, err)
	assert.Equal(t, "0.05", got.Rate.String())
	assert.Equal(t, "0.07", got.TotalRate.String())
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "weeth", params["vault_type"])
}

func TestGetVaultStatistics_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_vault_statistics", []any{
		map[string]any{
			"vault_name": "rswETHC", "block_number": int64(12345), "block_timestamp": int64(1700000000),
			"total_supply": "1000", "usd_tvl": "2700000", "usd_value": "2700",
			"base_value": "1.05", "underlying_value": "1.07", "subaccount_value_at_last_trade": "2650",
		},
	})
	got, err := api.GetVaultStatistics(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "rswETHC", got[0].VaultName)
	assert.Equal(t, "2700000", got[0].USDTVL.String())
}
