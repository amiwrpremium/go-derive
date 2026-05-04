package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestSubAccount_Decode(t *testing.T) {
	payload := `{
		"subaccount_id": 7,
		"owner_address": "0x1111111111111111111111111111111111111111",
		"margin_type": "PM",
		"is_under_liquidation": false,
		"subaccount_value": "100",
		"initial_margin": "50",
		"maintenance_margin": "30"
	}`
	var sa types.SubAccount
	require.NoError(t, json.Unmarshal([]byte(payload), &sa))
	assert.Equal(t, int64(7), sa.SubaccountID)
	assert.Equal(t, "PM", sa.MarginType)
	assert.False(t, sa.IsUnderLiquidation)
}

func TestSubAccount_RoundTrip(t *testing.T) {
	in := types.SubAccount{
		SubaccountID:      1,
		OwnerAddress:      types.MustAddress("0x1111111111111111111111111111111111111111"),
		MarginType:        "SM",
		SubaccountValue:   types.MustDecimal("0"),
		InitialMargin:     types.MustDecimal("0"),
		MaintenanceMargin: types.MustDecimal("0"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.SubAccount
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.SubaccountID, out.SubaccountID)
	assert.Equal(t, in.OwnerAddress.String(), out.OwnerAddress.String())
}
