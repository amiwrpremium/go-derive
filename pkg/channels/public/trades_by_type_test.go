package public_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/channels/public"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestTradesByType_Name(t *testing.T) {
	assert.Equal(t, "trades.perp.BTC",
		public.TradesByType{InstrumentType: enums.InstrumentTypePerp, Currency: "BTC"}.Name())
	assert.Equal(t, "trades.option.ETH",
		public.TradesByType{InstrumentType: enums.InstrumentTypeOption, Currency: "ETH"}.Name())
}

func TestTradesByType_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"trade_id":"t1","instrument_name":"BTC-PERP","direction":"buy","trade_price":"100","trade_amount":"1","timestamp":1700000000000}]`)
	v, err := public.TradesByType{}.Decode(raw)
	require.NoError(t, err)
	trades, ok := v.([]types.Trade)
	require.True(t, ok)
	require.Len(t, trades, 1)
	assert.Equal(t, "t1", trades[0].TradeID)
}

func TestTradesByType_Decode_Malformed(t *testing.T) {
	_, err := public.TradesByType{}.Decode([]byte(`{`))
	assert.Error(t, err)
}

func TestTradesByTypeTxStatus_Name(t *testing.T) {
	got := public.TradesByTypeTxStatus{
		InstrumentType: enums.InstrumentTypePerp,
		Currency:       "BTC",
		TxStatus:       enums.TxStatusSettled,
	}.Name()
	assert.Equal(t, "trades.perp.BTC.settled", got)
}

func TestTradesByTypeTxStatus_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"trade_id":"t1","instrument_name":"BTC-PERP","direction":"buy","trade_price":"100","trade_amount":"1","timestamp":1700000000000}]`)
	v, err := public.TradesByTypeTxStatus{TxStatus: enums.TxStatusSettled}.Decode(raw)
	require.NoError(t, err)
	trades, ok := v.([]types.Trade)
	require.True(t, ok)
	require.Len(t, trades, 1)
}

func TestTradesByTypeTxStatus_Decode_Malformed(t *testing.T) {
	_, err := public.TradesByTypeTxStatus{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
