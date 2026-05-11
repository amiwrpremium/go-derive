package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestGetFundingRateHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_funding_rate_history", map[string]any{
		"funding_rate_history": []any{
			map[string]any{"timestamp": int64(1700000000000), "funding_rate": "0.0001"},
		},
	})
	got, err := api.GetFundingRateHistory(context.Background(), types.FundingRateHistoryQuery{InstrumentName: "BTC-PERP"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "0.0001", got[0].FundingRate.String())
	assert.Equal(t, int64(1700000000000), got[0].Timestamp.Millis())
}

func TestGetSpotFeedHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_spot_feed_history", map[string]any{
		"currency": "BTC",
		"spot_feed_history": []any{
			map[string]any{"timestamp": int64(1700000000000), "timestamp_bucket": int64(1700000000000), "price": "50000"},
		},
	})
	currency, items, err := api.GetSpotFeedHistory(context.Background(), types.SpotFeedHistoryQuery{
		Currency:  "BTC",
		PeriodSec: 60,
	})
	require.NoError(t, err)
	assert.Equal(t, "BTC", currency)
	require.Len(t, items, 1)
	assert.Equal(t, "50000", items[0].Price.String())
}

func TestGetLatestSignedFeeds_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_latest_signed_feeds", map[string]any{
		"spot_data": map[string]any{
			"BTC": map[string]any{
				"currency": "BTC", "price": "50000", "confidence": "0.99",
				"timestamp": int64(1700000000000), "deadline": int64(1700000060000),
				"feed_source_type": "S",
				"signatures":       map[string]any{"signers": []any{"0xa"}, "signatures": []any{"0xs"}},
			},
		},
		"perp_data": map[string]any{},
		"fwd_data":  map[string]any{},
		"rate_data": map[string]any{},
		"vol_data":  map[string]any{},
	})
	feeds, err := api.GetLatestSignedFeeds(context.Background(), "", 0)
	require.NoError(t, err)
	require.NotNil(t, feeds)
	require.Contains(t, feeds.SpotData, "BTC")
	assert.Equal(t, "50000", feeds.SpotData["BTC"].Price.String())
}

func TestGetInterestRateHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_interest_rate_history", map[string]any{
		"interest_rates": []any{
			map[string]any{
				"block": int64(12345), "timestamp_sec": int64(1700000000),
				"borrow_apy": "0.08", "supply_apy": "0.04",
				"total_borrow": "5000000", "total_supply": "10000000",
			},
		},
		"pagination": map[string]any{"num_pages": 1, "count": 1},
	})
	rates, page, err := api.GetInterestRateHistory(context.Background(), types.InterestRateHistoryQuery{
		FromSec: 1700000000,
		ToSec:   1700100000,
	}, types.PageRequest{})
	require.NoError(t, err)
	require.Len(t, rates, 1)
	assert.Equal(t, "0.08", rates[0].BorrowAPY.String())
	assert.Equal(t, 1, page.Count)
}

func TestGetPerpImpactTWAP_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_perp_impact_twap", map[string]any{
		"currency":             "BTC",
		"mid_price_diff_twap":  "0.5",
		"ask_impact_diff_twap": "1.2",
		"bid_impact_diff_twap": "-0.8",
	})
	twap, err := api.GetPerpImpactTWAP(context.Background(), "BTC", 0, 1)
	require.NoError(t, err)
	assert.Equal(t, "BTC", twap.Currency)
	assert.Equal(t, "0.5", twap.MidPriceDiffTWAP.String())
	assert.Equal(t, "1.2", twap.AskImpactDiffTWAP.String())
	assert.Equal(t, "-0.8", twap.BidImpactDiffTWAP.String())
}
