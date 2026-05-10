package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestMarginWatch_Decode_PM(t *testing.T) {
	raw := []byte(`{
		"subaccount_id": 42,
		"currency": "USDC",
		"margin_type": "PM",
		"subaccount_value": "10000",
		"maintenance_margin": "-50.5",
		"valuation_timestamp": 1700000000
	}`)
	var m types.MarginWatch
	require.NoError(t, json.Unmarshal(raw, &m))
	assert.Equal(t, int64(42), m.SubaccountID)
	assert.Equal(t, enums.MarginTypePM, m.MarginType)
	assert.Equal(t, "10000", m.SubaccountValue.String())
	assert.Equal(t, "-50.5", m.MaintenanceMargin.String(),
		"negative maintenance margin signals the subaccount is below the liquidation floor")
	assert.Equal(t, int64(1700000000), m.ValuationTimestamp)
}

func TestMarginWatch_Decode_SM(t *testing.T) {
	raw := []byte(`{
		"subaccount_id": 7,
		"currency": "USDC",
		"margin_type": "SM",
		"subaccount_value": "100",
		"maintenance_margin": "10",
		"valuation_timestamp": 1700000060
	}`)
	var m types.MarginWatch
	require.NoError(t, json.Unmarshal(raw, &m))
	assert.Equal(t, enums.MarginTypeSM, m.MarginType)
}
