package public_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/channels/public"
)

func TestOrderBook_Name_DefaultsApplied(t *testing.T) {
	got := public.OrderBook{Instrument: "BTC-PERP"}.Name()
	assert.Equal(t, "orderbook.BTC-PERP.1.10", got)
}

func TestOrderBook_Name_ExplicitGroupAndDepth(t *testing.T) {
	got := public.OrderBook{Instrument: "ETH-PERP", Group: "10", Depth: 25}.Name()
	assert.Equal(t, "orderbook.ETH-PERP.10.25", got)
}

func TestOrderBook_Name_ExplicitGroupDefaultDepth(t *testing.T) {
	got := public.OrderBook{Instrument: "X", Group: "5", Depth: 0}.Name()
	assert.Equal(t, "orderbook.X.5.10", got)
	assert.Equal(t, "orderbook.X.5.10", public.OrderBook{Instrument: "X", Group: "5"}.Name())
}

func TestOrderBook_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`{
		"instrument_name":"BTC-PERP",
		"bids":[["100","1"]],
		"asks":[["101","2"]],
		"timestamp":1700000000000
	}`)
	v, err := public.OrderBook{}.Decode(raw)
	require.NoError(t, err)
	ob, ok := v.(derive.OrderBook)
	require.True(t, ok)
	assert.Equal(t, "BTC-PERP", ob.InstrumentName)
}

func TestOrderBook_Decode_Malformed(t *testing.T) {
	_, err := public.OrderBook{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
