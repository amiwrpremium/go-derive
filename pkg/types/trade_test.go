package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestTrade_Decode(t *testing.T) {
	payload := `{
		"trade_id": "T1",
		"order_id": "O1",
		"subaccount_id": 1,
		"instrument_name": "BTC-PERP",
		"direction": "buy",
		"trade_price": "65000",
		"trade_amount": "0.1",
		"mark_price": "65000",
		"index_price": "64999",
		"trade_fee": "0.5",
		"liquidity_role": "taker",
		"realized_pnl": "10",
		"timestamp": 1700000000000
	}`
	var tr types.Trade
	require.NoError(t, json.Unmarshal([]byte(payload), &tr))
	assert.Equal(t, "T1", tr.TradeID)
	assert.Equal(t, enums.DirectionBuy, tr.Direction)
	assert.Equal(t, enums.LiquidityRoleTaker, tr.LiquidityRole)
}

func TestTrade_OmitsEmpty(t *testing.T) {
	in := types.Trade{
		TradeID:        "T1",
		InstrumentName: "BTC-PERP",
		Direction:      enums.DirectionBuy,
		TradePrice:     types.MustDecimal("100"),
		TradeAmount:    types.MustDecimal("1"),
		MarkPrice:      types.MustDecimal("100"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	s := string(b)
	assert.NotContains(t, s, `"order_id"`)
	assert.NotContains(t, s, `"liquidity_role"`)
}
