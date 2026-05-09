package private_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/channels/private"
)

func TestPrivateTrades_Name(t *testing.T) {
	assert.Equal(t, "subaccount.7.trades", private.Trades{SubaccountID: 7}.Name())
}

func TestPrivateTrades_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"trade_id":"T","instrument_name":"BTC-PERP","direction":"buy","trade_price":"1","trade_amount":"1","mark_price":"1","timestamp":1700000000000}]`)
	v, err := private.Trades{}.Decode(raw)
	require.NoError(t, err)
	trades, ok := v.([]derive.Trade)
	require.True(t, ok)
	require.Len(t, trades, 1)
}

func TestPrivateTrades_Decode_Malformed(t *testing.T) {
	_, err := private.Trades{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
