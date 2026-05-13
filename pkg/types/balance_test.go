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
		"amount_step": "0.000001",
		"mark_price": "1",
		"mark_value": "10000",
		"average_price": "1",
		"average_price_excl_fees": "1",
		"cumulative_interest": "5",
		"pending_interest": "0.1",
		"initial_margin": "100",
		"maintenance_margin": "50",
		"open_orders_margin": "20",
		"delta": "10000",
		"delta_currency": "USDC",
		"realized_pnl": "5",
		"realized_pnl_excl_fees": "6",
		"unrealized_pnl": "1",
		"unrealized_pnl_excl_fees": "2",
		"total_fees": "3",
		"creation_timestamp": 1700000000000
	}`
	var c types.Collateral
	require.NoError(t, json.Unmarshal([]byte(payload), &c))
	assert.Equal(t, "USDC", c.AssetName)
	assert.Equal(t, enums.AssetTypeERC20, c.AssetType)
	assert.Equal(t, "10000", c.Amount.String())
	assert.Equal(t, "USDC", c.DeltaCurrency)
	assert.Equal(t, "20", c.OpenOrdersMargin.String())
	assert.Equal(t, "6", c.RealizedPNLExclFees.String())
	assert.Equal(t, "2", c.UnrealizedPNLExclFees.String())
	assert.Equal(t, "3", c.TotalFees.String())
	assert.Equal(t, int64(1700000000000), c.CreationTimestamp.Millis())
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
