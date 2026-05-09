package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFundingRateHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_funding_rate_history", map[string]any{
		"funding_rate_history": []any{
			map[string]any{"timestamp": int64(1700000000000), "funding_rate": "0.0001"},
		},
	})
	got, err := api.GetFundingRateHistory(context.Background(), map[string]any{"instrument_name": "BTC-PERP"})
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
	currency, items, err := api.GetSpotFeedHistory(context.Background(), map[string]any{
		"currency": "BTC", "period": int64(60),
		"start_timestamp": 0, "end_timestamp": 1,
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
	feeds, err := api.GetLatestSignedFeeds(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, feeds)
	require.Contains(t, feeds.SpotData, "BTC")
	assert.Equal(t, "50000", feeds.SpotData["BTC"].Price.String())
}

func TestGetPerpImpactTWAP_Raw_Reachable(t *testing.T) {
	// `public/get_perp_impact_twap` is the one feed method that stays
	// raw because it isn't published in the v2.2 OAS. The wrapper still
	// works — verify it forwards params and returns the bytes.
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_perp_impact_twap", map[string]any{"impact_price": "100"})
	raw, err := api.GetPerpImpactTWAP(context.Background(), map[string]any{
		"currency": "BTC", "start_time": 0, "end_time": 1,
	})
	require.NoError(t, err)
	assert.JSONEq(t, `{"impact_price":"100"}`, string(raw))
}
