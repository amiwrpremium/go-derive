package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestGetIndexChartData_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_index_chart_data", []any{
		map[string]any{
			"timestamp": int64(1700000000000), "timestamp_bucket": int64(1700000000000),
			"price": "2500", "open_price": "2495", "high_price": "2510",
			"low_price": "2490", "close_price": "2500",
		},
	})
	got, err := api.GetIndexChartData(context.Background(), types.IndexChartQuery{
		Currency:  "ETH",
		PeriodSec: 60,
	})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "2500", got[0].Price.String())
}

func TestGetTradingViewChartData_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_tradingview_chart_data", []any{
		map[string]any{
			"timestamp": int64(1700000000000), "timestamp_bucket": int64(1700000000000),
			"open_price": "65000", "high_price": "65100", "low_price": "64900", "close_price": "65050",
			"volume_contracts": "10", "volume_usd": "650500",
		},
	})
	got, err := api.GetTradingViewChartData(context.Background(), types.TradingViewChartQuery{
		InstrumentName: "BTC-PERP",
		PeriodSec:      60,
	})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "65000", got[0].OpenPrice.String())
	assert.Equal(t, "650500", got[0].VolumeUSD.String())
}

func TestGetSpotFeedHistoryCandles_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_spot_feed_history_candles", map[string]any{
		"currency": "BTC",
		"spot_feed_history": []any{
			map[string]any{
				"timestamp": int64(1700000000000), "timestamp_bucket": int64(1700000000000),
				"price": "65000", "open_price": "65000", "high_price": "65100",
				"low_price": "64900", "close_price": "65050",
			},
		},
	})
	currency, candles, err := api.GetSpotFeedHistoryCandles(context.Background(), types.SpotFeedHistoryCandlesQuery{
		Currency:  "BTC",
		PeriodSec: 60,
	})
	require.NoError(t, err)
	assert.Equal(t, "BTC", currency)
	require.Len(t, candles, 1)
	assert.Equal(t, "65050", candles[0].ClosePrice.String())
}
