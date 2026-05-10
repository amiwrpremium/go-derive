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

func TestTicker_Name_DefaultInterval(t *testing.T) {
	got := public.Ticker{Instrument: "BTC-PERP"}.Name()
	assert.Equal(t, "ticker.BTC-PERP.1000", got)
}

func TestTicker_Name_ExplicitInterval(t *testing.T) {
	got := public.Ticker{Instrument: "BTC-PERP", Interval: "100"}.Name()
	assert.Equal(t, "ticker.BTC-PERP.100", got)
}

func TestTicker_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`{
		"timestamp": 1700000000000,
		"instrument_ticker": {
			"instrument_name": "BTC-PERP",
			"instrument_type": "perp",
			"is_active": true,
			"base_currency": "BTC",
			"quote_currency": "USD",
			"base_asset_address": "0x1111111111111111111111111111111111111111",
			"base_asset_sub_id": "0",
			"amount_step": "0.001", "minimum_amount": "0.001", "maximum_amount": "100",
			"tick_size": "0.5",
			"base_fee": "0.1", "maker_fee_rate": "0.0003", "taker_fee_rate": "0.0005",
			"mark_price_fee_rate_cap": null,
			"scheduled_activation": 1700000000, "scheduled_deactivation": 9223372036854775807,
			"best_bid_price": "65000", "best_bid_amount": "1.5",
			"best_ask_price": "65010", "best_ask_amount": "2.0",
			"mark_price": "65005", "index_price": "65003",
			"min_price": "64000", "max_price": "66000",
			"timestamp": 1700000000000,
			"option_details": null, "perp_details": {"index":"BTC"},
			"option_pricing": null,
			"stats": {"contract_volume":"1","high":"1","low":"1","num_trades":"1","open_interest":"1","percent_change":"0","usd_change":"0"}
		}
	}`)
	v, err := public.Ticker{}.Decode(raw)
	require.NoError(t, err)
	feed, ok := v.(types.InstrumentTickerFeed)
	require.True(t, ok, "Decode must return types.InstrumentTickerFeed")
	assert.Equal(t, "BTC-PERP", feed.Ticker.InstrumentName)
	assert.Equal(t, "65005", feed.Ticker.MarkPrice.String())
	assert.Equal(t, int64(1700000000000), feed.Timestamp.Millis())
}

func TestTicker_Decode_Malformed(t *testing.T) {
	_, err := public.Ticker{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
