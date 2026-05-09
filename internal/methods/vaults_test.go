package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVaultBalances_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_vault_balances", []any{
		map[string]any{
			"name": "rswETHC", "address": "0x1111111111111111111111111111111111111111",
			"chain_id": int64(1), "vault_asset_type": "rswETHC", "amount": "12.5",
		},
	})
	got, err := api.GetVaultBalances(context.Background(), map[string]any{
		"wallet": "0x2222222222222222222222222222222222222222",
	})
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
	shares, page, err := api.GetVaultShare(context.Background(), map[string]any{
		"vault_name": "rswETHC", "from_timestamp_sec": int64(1700000000), "to_timestamp_sec": int64(1700100000),
	})
	require.NoError(t, err)
	require.Len(t, shares, 1)
	assert.Equal(t, "1.07", shares[0].UnderlyingValue.String())
	assert.Equal(t, 1, page.Count)
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
