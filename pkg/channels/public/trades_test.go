package public_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/channels/public"
)

func TestTrades_Name(t *testing.T) {
	assert.Equal(t, "trades.BTC-PERP", public.Trades{Instrument: "BTC-PERP"}.Name())
	assert.Equal(t, "trades.", public.Trades{}.Name())
}

func TestTrades_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"trade_id":"T1","instrument_name":"BTC-PERP","direction":"buy","trade_price":"65000","trade_amount":"0.1","mark_price":"65000","timestamp":1700000000000}]`)
	v, err := public.Trades{}.Decode(raw)
	require.NoError(t, err)
	trades, ok := v.([]derive.Trade)
	require.True(t, ok)
	require.Len(t, trades, 1)
	assert.Equal(t, "T1", trades[0].TradeID)
}

func TestTrades_Decode_EmptyArray(t *testing.T) {
	v, err := public.Trades{}.Decode(json.RawMessage(`[]`))
	require.NoError(t, err)
	trades := v.([]derive.Trade)
	assert.Empty(t, trades)
}

func TestTrades_Decode_Malformed(t *testing.T) {
	_, err := public.Trades{}.Decode([]byte(`not-json`))
	assert.Error(t, err)
}
