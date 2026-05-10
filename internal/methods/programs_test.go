package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMakerPrograms_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_maker_programs", []any{
		map[string]any{
			"name": "options-q1", "asset_types": []any{"option"}, "currencies": []any{"ETH"},
			"min_notional": "10000", "rewards": map[string]any{"DRV": "1000"},
			"start_timestamp": int64(1700000000), "end_timestamp": int64(1702592000),
		},
	})
	got, err := api.GetMakerPrograms(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "options-q1", got[0].Name)
}

func TestGetMakerProgramScores_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_maker_program_scores", map[string]any{
		"program": map[string]any{
			"name": "options-q1", "asset_types": []any{"option"}, "currencies": []any{"ETH"},
			"min_notional": "10000", "rewards": map[string]any{"DRV": "1000"},
			"start_timestamp": int64(1700000000), "end_timestamp": int64(1702592000),
		},
		"scores": []any{
			map[string]any{
				"wallet":         "0x1111111111111111111111111111111111111111",
				"coverage_score": "0.8", "quality_score": "0.9", "holder_boost": "1",
				"volume": "1000", "volume_multiplier": "1", "total_score": "0.72",
			},
		},
		"total_score": "1", "total_volume": "1000",
	})
	got, err := api.GetMakerProgramScores(context.Background(), map[string]any{
		"program_name": "options-q1", "epoch_start_timestamp": int64(1700000000),
	})
	require.NoError(t, err)
	assert.Equal(t, "options-q1", got.Program.Name)
	require.Len(t, got.Scores, 1)
	assert.Equal(t, "0.72", got.Scores[0].TotalScore.String())
}

func TestGetReferralPerformance_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_referral_performance", map[string]any{
		"referral_code":               "ALICE",
		"fee_share_percentage":        "0.2",
		"stdrv_balance":               "1000",
		"total_notional_volume":       "500000",
		"total_referred_fees":         "1500",
		"total_fee_rewards":           "300",
		"total_builder_fee_collected": "50",
		"rewards":                     map[string]any{},
	})
	got, err := api.GetReferralPerformance(context.Background(), map[string]any{
		"start_ms": int64(1700000000000), "end_ms": int64(1700100000000),
	})
	require.NoError(t, err)
	assert.Equal(t, "ALICE", got.ReferralCode)
	assert.Equal(t, "0.2", got.FeeSharePercentage.String())
}

func TestGetDetailedMakerSnapshotHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_detailed_maker_snapshot_history", map[string]any{
		"program": map[string]any{
			"name": "options-q1", "asset_types": []any{"option"}, "currencies": []any{"ETH"},
			"min_notional": "10000", "rewards": map[string]any{"DRV": "100"},
			"start_timestamp": int64(1700000000), "end_timestamp": int64(1702592000),
		},
		"snapshots": []any{
			map[string]any{
				"wallet":                "0x1111111111111111111111111111111111111111",
				"asset_type":            "option",
				"currency":              "ETH",
				"instrument_name":       "ETH-25DEC25-3000-C",
				"is_bid":                true,
				"timestamp":             int64(1700000000000),
				"mid_price":             "100",
				"best_price":            "99.5",
				"index_price":           "2500",
				"notional":              "10000",
				"scaled_notional":       "9500",
				"total_scaled_notional": "1000000",
			},
		},
		"pagination": map[string]any{"num_pages": 1, "count": 1},
	})
	got, err := api.GetDetailedMakerSnapshotHistory(context.Background(), map[string]any{
		"program_name":          "options-q1",
		"epoch_start_timestamp": int64(1700000000),
	})
	require.NoError(t, err)
	assert.Equal(t, "options-q1", got.Program.Name)
	require.Len(t, got.Snapshots, 1)
	assert.Equal(t, "ETH-25DEC25-3000-C", got.Snapshots[0].InstrumentName)
}
