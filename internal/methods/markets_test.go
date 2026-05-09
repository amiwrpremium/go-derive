package methods_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestGetInstruments(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_instruments", []map[string]any{
		{"instrument_name": "BTC-PERP", "instrument_type": "perp", "is_active": true,
			"base_currency": "BTC", "quote_currency": "USDC",
			"tick_size": "0.5", "minimum_amount": "0.001", "maximum_amount": "100", "amount_step": "0.001"},
	})

	insts, err := api.GetInstruments(context.Background(), "BTC", enums.InstrumentTypePerp)
	require.NoError(t, err)
	require.Len(t, insts, 1)
	assert.Equal(t, "BTC-PERP", insts[0].Name)

	last := ft.LastCall()
	assert.Equal(t, "public/get_instruments", last.Method)
	params := paramsAsMap(t, last.Params)
	assert.Equal(t, "BTC", params["currency"])
	assert.Equal(t, "perp", params["instrument_type"])
	assert.Equal(t, false, params["expired"])
}

func TestGetInstruments_NoFilters(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_instruments", []map[string]any{})

	_, err := api.GetInstruments(context.Background(), "", "")
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	_, hasCurrency := params["currency"]
	_, hasKind := params["instrument_type"]
	assert.False(t, hasCurrency)
	assert.False(t, hasKind)
}

func TestGetInstrument(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_instrument", map[string]any{
		"instrument_name": "BTC-PERP",
		"instrument_type": "perp",
		"is_active":       true,
		"base_currency":   "BTC",
		"quote_currency":  "USDC",
		"tick_size":       "0.5",
		"minimum_amount":  "0.001",
		"maximum_amount":  "100",
		"amount_step":     "0.001",
	})

	got, err := api.GetInstrument(context.Background(), "BTC-PERP")
	require.NoError(t, err)
	assert.Equal(t, "BTC-PERP", got.Name)
}

func TestGetTicker(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_ticker", map[string]any{
		"instrument_name": "BTC-PERP",
		"best_bid_price":  "100", "best_bid_amount": "1",
		"best_ask_price": "101", "best_ask_amount": "2",
		"mark_price":  "100.5",
		"index_price": "100.5",
		"timestamp":   1700000000000,
	})
	got, err := api.GetTicker(context.Background(), "BTC-PERP")
	require.NoError(t, err)
	assert.Equal(t, "BTC-PERP", got.InstrumentName)
}

func TestGetPublicTradeHistory(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_trade_history", map[string]any{
		"trades":     []any{},
		"pagination": map[string]any{"num_pages": 1, "count": 0, "current_page": 1, "page_size": 50},
	})
	trades, page, err := api.GetPublicTradeHistory(context.Background(), "BTC-PERP", types.PageRequest{Page: 2, PageSize: 50})
	require.NoError(t, err)
	assert.Empty(t, trades)
	assert.Equal(t, 1, page.NumPages)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, float64(2), params["page"])
	assert.Equal(t, float64(50), params["page_size"])
}

func TestGetTime(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_time", 1700000000000)
	got, err := api.GetTime(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(1700000000000), got)
}

func TestGetCurrencies(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_all_currencies", []map[string]any{
		{"currency": "USDC"}, {"currency": "WETH"},
	})
	got, err := api.GetCurrencies(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"USDC", "WETH"}, got)
}

func TestPublicMethods_PropagateUnhandledTransportError(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	// No handler registered → FakeTransport returns "unhandled" error.
	_, err := api.GetTime(context.Background())
	require.Error(t, err)
}

func TestGetStatistics_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/statistics", map[string]any{
		"daily_fees":            "100",
		"daily_notional_volume": "1000000",
		"daily_premium_volume":  "50000",
		"daily_trades":          int64(250),
		"open_interest":         "500",
		"total_fees":            "10000",
		"total_notional_volume": "100000000",
		"total_premium_volume":  "500000",
		"total_trades":          int64(25000),
	})
	got, err := api.GetStatistics(context.Background(), "BTC-PERP")
	require.NoError(t, err)
	assert.Equal(t, int64(250), got.DailyTrades)
	assert.Equal(t, "500", got.OpenInterest.String())
}

// silence unused json import warning on platforms where compiler is strict.
var _ = json.Marshal
