package private_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/channels/private"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestPrivateTrades_Name(t *testing.T) {
	assert.Equal(t, "subaccount.7.trades", private.Trades{SubaccountID: 7}.Name())
}

func TestPrivateTrades_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"trade_id":"T","instrument_name":"BTC-PERP","direction":"buy","trade_price":"1","trade_amount":"1","mark_price":"1","timestamp":1700000000000}]`)
	v, err := private.Trades{}.Decode(raw)
	require.NoError(t, err)
	trades, ok := v.([]types.Trade)
	require.True(t, ok)
	require.Len(t, trades, 1)
}

func TestPrivateTrades_Decode_Malformed(t *testing.T) {
	_, err := private.Trades{}.Decode([]byte(`{`))
	assert.Error(t, err)
}

func TestPrivateTradesByTxStatus_Name(t *testing.T) {
	got := private.TradesByTxStatus{SubaccountID: 7, TxStatus: enums.TxStatusSettled}.Name()
	assert.Equal(t, "subaccount.7.trades.settled", got)
}

func TestPrivateTradesByTxStatus_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"trade_id":"T","instrument_name":"BTC-PERP","direction":"buy","trade_price":"1","trade_amount":"1","mark_price":"1","timestamp":1700000000000}]`)
	v, err := private.TradesByTxStatus{TxStatus: enums.TxStatusSettled}.Decode(raw)
	require.NoError(t, err)
	trades, ok := v.([]types.Trade)
	require.True(t, ok)
	require.Len(t, trades, 1)
}

func TestPrivateTradesByTxStatus_Decode_Malformed(t *testing.T) {
	_, err := private.TradesByTxStatus{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
