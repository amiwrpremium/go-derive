// Package methods is the shared implementation of every JSON-RPC method
// Derive exposes. Both pkg/rest.Client and pkg/ws.Client embed *API so that
// each method is defined exactly once, parameterised by the underlying
// transport.
//
// Public methods are unauthenticated; private methods require Signer to be
// non-nil. Private methods that mutate orders also use the Domain to sign
// the per-action EIP-712 hash.
package methods

import (
	"context"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// GetIndexChartData returns OHLC candles for one currency's index
// price feed over the requested window. Public.
//
// PeriodSec accepts the engine's bucket sizes (60, 300, 900, 1800,
// 3600, 14400, 28800, 86400, 604800).
//
// Each bar shares its shape with [API.GetSpotFeedHistoryCandles];
// the difference is the data source (index feed vs spot feed).
func (a *API) GetIndexChartData(ctx context.Context, q types.IndexChartQuery) ([]types.SpotFeedCandle, error) {
	params := map[string]any{
		"currency": q.Currency,
		"period":   q.PeriodSec,
	}
	addHistoryWindow(params, q.HistoryWindow)
	var resp []types.SpotFeedCandle
	if err := a.call(ctx, "public/get_index_chart_data", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTradingViewChartData returns TradingView UDF-format OHLC bars
// for one instrument over the requested window. Public.
//
// Distinct from [API.GetIndexChartData] in that this method's bars
// carry per-bar volume in both contracts and USD.
func (a *API) GetTradingViewChartData(ctx context.Context, q types.TradingViewChartQuery) ([]types.TradingViewChart, error) {
	params := map[string]any{
		"instrument_name": q.InstrumentName,
		"period":          q.PeriodSec,
	}
	addHistoryWindow(params, q.HistoryWindow)
	var resp []types.TradingViewChart
	if err := a.call(ctx, "public/get_tradingview_chart_data", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetSpotFeedHistoryCandles returns OHLC candles for one currency's
// spot feed over the requested window. Public.
//
// Returns the currency the response is keyed against alongside the
// per-bucket samples (the wire response repeats `currency` next to
// the candle list).
func (a *API) GetSpotFeedHistoryCandles(ctx context.Context, q types.SpotFeedHistoryCandlesQuery) (currency string, candles []types.SpotFeedCandle, err error) {
	params := map[string]any{
		"currency": q.Currency,
		"period":   q.PeriodSec,
	}
	addHistoryWindow(params, q.HistoryWindow)
	var resp struct {
		Currency        string                 `json:"currency"`
		SpotFeedHistory []types.SpotFeedCandle `json:"spot_feed_history"`
	}
	if err := a.call(ctx, "public/get_spot_feed_history_candles", params, &resp); err != nil {
		return "", nil, err
	}
	return resp.Currency, resp.SpotFeedHistory, nil
}
