package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestPosition_DecodeFull(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"instrument_type": "perp",
		"subaccount_id": 42,
		"amount": "0.5",
		"amount_step": "0.001",
		"average_price": "65000",
		"average_price_excl_fees": "64980",
		"mark_price": "65500",
		"mark_value": "32750",
		"index_price": "65500",
		"leverage": "5",
		"liquidation_price": "10000",
		"unrealized_pnl": "250",
		"unrealized_pnl_excl_fees": "260",
		"realized_pnl": "10",
		"realized_pnl_excl_fees": "11",
		"total_fees": "2",
		"open_orders_margin": "100",
		"cumulative_funding": "1",
		"pending_funding": "0.1"
	}`
	var p types.Position
	require.NoError(t, json.Unmarshal([]byte(payload), &p))
	assert.Equal(t, "BTC-PERP", p.InstrumentName)
	assert.Equal(t, enums.InstrumentTypePerp, p.InstrumentType)
	assert.Equal(t, int64(42), p.SubaccountID)
	assert.Equal(t, "0.5", p.Amount.String())
	assert.Equal(t, "0.001", p.AmountStep.String())
	assert.Equal(t, "64980", p.AveragePriceExclFees.String())
	assert.Equal(t, "260", p.UnrealizedPNLExclFees.String())
	assert.Equal(t, "11", p.RealizedPNLExclFees.String())
	assert.Equal(t, "2", p.TotalFees.String())
	assert.Equal(t, "5", p.Leverage.String())
}

func TestPosition_DecodeMinimal(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"instrument_type": "perp",
		"amount": "0",
		"average_price": "0",
		"mark_price": "0",
		"mark_value": "0",
		"unrealized_pnl": "0",
		"realized_pnl": "0"
	}`
	var p types.Position
	require.NoError(t, json.Unmarshal([]byte(payload), &p))
	assert.True(t, p.Amount.IsZero())
}

func TestPosition_RoundTrip(t *testing.T) {
	in := types.Position{
		InstrumentName: "BTC-PERP",
		InstrumentType: enums.InstrumentTypePerp,
		Amount:         types.MustDecimal("1"),
		AveragePrice:   types.MustDecimal("65000"),
		MarkPrice:      types.MustDecimal("65500"),
		MarkValue:      types.MustDecimal("65500"),
		UnrealizedPNL:  types.MustDecimal("500"),
		RealizedPNL:    types.MustDecimal("0"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.Position
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.InstrumentName, out.InstrumentName)
	assert.Equal(t, in.UnrealizedPNL.String(), out.UnrealizedPNL.String())
}
