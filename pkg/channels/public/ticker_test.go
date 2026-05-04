package public_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/channels/public"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestTickerSlim_Name_DefaultInterval(t *testing.T) {
	got := public.TickerSlim{Instrument: "BTC-PERP"}.Name()
	assert.Equal(t, "ticker_slim.BTC-PERP.1000", got)
}

func TestTickerSlim_Name_ExplicitInterval(t *testing.T) {
	got := public.TickerSlim{Instrument: "BTC-PERP", Interval: "100"}.Name()
	assert.Equal(t, "ticker_slim.BTC-PERP.100", got)
}

func TestTickerSlim_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`{
		"timestamp": 1700000000000,
		"instrument_ticker": {
			"t": 1700000000000,
			"A": "0.5", "a": "78758.5",
			"B": "0.4", "b": "78752.1",
			"M": "78755", "I": "78760",
			"f": "0.0001"
		}
	}`)
	v, err := public.TickerSlim{}.Decode(raw)
	require.NoError(t, err)
	tk, ok := v.(types.TickerSlim)
	require.True(t, ok)
	assert.Equal(t, "78752.1", tk.Ticker.BestBidPrice.String())
	assert.Equal(t, "78758.5", tk.Ticker.BestAskPrice.String())
	assert.Equal(t, "78755", tk.Ticker.MarkPrice.String())
}

func TestTickerSlim_Decode_Malformed(t *testing.T) {
	_, err := public.TickerSlim{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
