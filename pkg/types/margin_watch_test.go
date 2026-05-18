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

func TestMarginWatch_Decode_WithNestedBreakdown(t *testing.T) {
	raw := []byte(`{
		"subaccount_id": 42,
		"currency": "USDC",
		"margin_type": "PM",
		"subaccount_value": "10000",
		"initial_margin": "-100",
		"maintenance_margin": "-50.5",
		"valuation_timestamp": 1700000000,
		"collaterals": [
			{
				"asset_name": "USDC",
				"asset_type": "erc20",
				"amount": "9500",
				"mark_price": "1",
				"mark_value": "9500",
				"delta": "0",
				"delta_currency": "USDC",
				"initial_margin": "9500",
				"maintenance_margin": "9500"
			}
		],
		"positions": [
			{
				"instrument_name": "ETH-PERP",
				"instrument_type": "perp",
				"amount": "-5",
				"delta": "-5",
				"gamma": "0",
				"theta": "0",
				"vega": "0",
				"index_price": "2500",
				"mark_price": "2510",
				"mark_value": "-12550",
				"initial_margin": "-9600",
				"maintenance_margin": "-9550",
				"liquidation_price": "2700"
			}
		]
	}`)
	var m types.MarginWatch
	require.NoError(t, json.Unmarshal(raw, &m))
	assert.Equal(t, "-100", m.InitialMargin.String())
	require.Len(t, m.Collaterals, 1)
	assert.Equal(t, "USDC", m.Collaterals[0].AssetName)
	assert.Equal(t, "9500", m.Collaterals[0].MarkValue.String())
	require.Len(t, m.Positions, 1)
	assert.Equal(t, "ETH-PERP", m.Positions[0].InstrumentName)
	assert.Equal(t, "2700", m.Positions[0].LiquidationPrice.String())
	assert.Equal(t, "-5", m.Positions[0].Amount.String())
}
