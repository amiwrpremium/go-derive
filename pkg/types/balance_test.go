package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestCollateral_Decode(t *testing.T) {
	payload := `{
		"asset_name": "USDC",
		"asset_type": "erc20",
		"amount": "10000",
		"mark_price": "1",
		"mark_value": "10000",
		"cumulative_interest": "5",
		"pending_interest": "0.1",
		"initial_margin": "100",
		"maintenance_margin": "50"
	}`
	var c types.Collateral
	require.NoError(t, json.Unmarshal([]byte(payload), &c))
	assert.Equal(t, "USDC", c.AssetName)
	assert.Equal(t, enums.AssetTypeERC20, c.AssetType)
	assert.Equal(t, "10000", c.Amount.String())
}

func TestCollateral_RoundTrip(t *testing.T) {
	in := types.Collateral{
		AssetName: "USDC", AssetType: enums.AssetTypeERC20,
		Amount:    types.MustDecimal("100"),
		MarkValue: types.MustDecimal("100"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.Collateral
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.AssetName, out.AssetName)
	assert.Equal(t, in.Amount.String(), out.Amount.String())
}

func TestBalance_Decode(t *testing.T) {
	payload := `{
		"subaccount_id": 123,
		"subaccount_value": "10000",
		"initial_margin": "5000",
		"maintenance_margin": "3000",
		"collaterals": [{"asset_name": "USDC", "asset_type": "erc20", "amount": "10000", "mark_value": "10000"}],
		"positions": []
	}`
	var b types.Balance
	require.NoError(t, json.Unmarshal([]byte(payload), &b))
	assert.Equal(t, int64(123), b.SubaccountID)
	require.Len(t, b.Collaterals, 1)
	assert.Equal(t, "USDC", b.Collaterals[0].AssetName)
	assert.Empty(t, b.Positions)
}

func TestBalance_OmitsEmptyPositionsOnMarshal(t *testing.T) {
	in := types.Balance{
		SubaccountID:      1,
		SubaccountValue:   types.MustDecimal("0"),
		InitialMargin:     types.MustDecimal("0"),
		MaintenanceMargin: types.MustDecimal("0"),
		Collaterals:       []types.Collateral{},
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.NotContains(t, string(b), "positions")
}

func TestBalanceUpdate_Decode(t *testing.T) {
	payload := `[
		{"name":"USDC","new_balance":"10500","previous_balance":"10000","update_type":"trade"},
		{"name":"BTC-PERP","new_balance":"-0.5","previous_balance":"0","update_type":"trade"}
	]`
	var got []types.BalanceUpdate
	require.NoError(t, json.Unmarshal([]byte(payload), &got))
	require.Len(t, got, 2)
	assert.Equal(t, "USDC", got[0].Name)
	assert.Equal(t, "10500", got[0].NewBalance.String())
	assert.Equal(t, "10000", got[0].PreviousBalance.String())
	assert.Equal(t, enums.BalanceUpdateTrade, got[0].UpdateType)
	assert.Equal(t, "BTC-PERP", got[1].Name)
	assert.Equal(t, "-0.5", got[1].NewBalance.String())
}
