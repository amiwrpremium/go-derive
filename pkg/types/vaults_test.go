package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestVaultBalance_Decode(t *testing.T) {
	raw := []byte(`{
		"name": "rswETHC",
		"address": "0x1111111111111111111111111111111111111111",
		"chain_id": 1,
		"vault_asset_type": "rswETHC",
		"amount": "12.5"
	}`)
	var v types.VaultBalance
	require.NoError(t, json.Unmarshal(raw, &v))
	assert.Equal(t, "rswETHC", v.Name)
	assert.Equal(t, "0x1111111111111111111111111111111111111111", v.Address.String())
	assert.Equal(t, int64(1), v.ChainID)
	assert.Equal(t, "12.5", v.Amount.String())
}

func TestVaultShare_Decode(t *testing.T) {
	raw := []byte(`{
		"block_number": 12345,
		"block_timestamp": 1700000000,
		"base_value": "1.05",
		"underlying_value": "1.07",
		"usd_value": "2700"
	}`)
	var s types.VaultShare
	require.NoError(t, json.Unmarshal(raw, &s))
	assert.Equal(t, int64(12345), s.BlockNumber)
	assert.Equal(t, "1.07", s.UnderlyingValue.String())
	assert.Equal(t, "2700", s.USDValue.String())
}

func TestVaultShare_NullUnderlying(t *testing.T) {
	raw := []byte(`{
		"block_number": 12345,
		"block_timestamp": 1700000000,
		"base_value": "1.05",
		"underlying_value": null,
		"usd_value": "2700"
	}`)
	var s types.VaultShare
	require.NoError(t, json.Unmarshal(raw, &s))
	assert.Equal(t, "0", s.UnderlyingValue.String())
}

func TestVaultStatistics_Decode(t *testing.T) {
	raw := []byte(`{
		"vault_name": "rswETHC",
		"block_number": 12345,
		"block_timestamp": 1700000000,
		"total_supply": "1000",
		"usd_tvl": "2700000",
		"usd_value": "2700",
		"base_value": "1.05",
		"underlying_value": "1.07",
		"subaccount_value_at_last_trade": "2650"
	}`)
	var v types.VaultStatistics
	require.NoError(t, json.Unmarshal(raw, &v))
	assert.Equal(t, "rswETHC", v.VaultName)
	assert.Equal(t, "1000", v.TotalSupply.String())
	assert.Equal(t, "2700000", v.USDTVL.String())
	assert.Equal(t, "2650", v.SubaccountValueAtLastTrade.String())
}

func TestVaultStatistics_NullablesAbsent(t *testing.T) {
	// Pre-trade vaults: subaccount_value_at_last_trade is null;
	// vaults without an underlying: underlying_value is null. Both
	// must decode to zero-value Decimals.
	raw := []byte(`{
		"vault_name": "rswETHC",
		"block_number": 12345,
		"block_timestamp": 1700000000,
		"total_supply": "0",
		"usd_tvl": "0",
		"usd_value": "0",
		"base_value": "1",
		"underlying_value": null,
		"subaccount_value_at_last_trade": null
	}`)
	var v types.VaultStatistics
	require.NoError(t, json.Unmarshal(raw, &v))
	assert.Equal(t, "0", v.UnderlyingValue.String())
	assert.Equal(t, "0", v.SubaccountValueAtLastTrade.String())
}
