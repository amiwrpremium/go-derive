package methods_test

import (
	"context"
	"testing"

	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllPoints_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_all_points", map[string]any{
		"total_notional_volume": "100000",
		"total_users":           int64(50),
		"points":                map[string]any{"0xabc": "100", "0xdef": "200"},
	})
	got, err := api.GetAllPoints(context.Background(), types.AllPointsQuery{Program: "trading"})
	require.NoError(t, err)
	assert.Equal(t, int64(50), got.TotalUsers)
	assert.NotEmpty(t, got.Points)
}

func TestGetPoints_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_points", map[string]any{
		"flag":                    "",
		"parent":                  "",
		"percent_share_of_points": "0.01",
		"total_notional_volume":   "1000",
		"total_users":             int64(50),
		"user_rank":               int64(7),
		"points":                  map[string]any{"trading": "100"},
	})
	got, err := api.GetPoints(context.Background(), types.PointsQuery{Program: "trading", Wallet: "0xabc"})
	require.NoError(t, err)
	assert.Equal(t, "1000", got.TotalNotionalVolume.String())
	assert.NotEmpty(t, got.Points)
}

func TestGetPointsLeaderboard_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_points_leaderboard", map[string]any{
		"pages":       int64(2),
		"total_users": int64(150),
		"leaderboard": []any{
			map[string]any{
				"rank": int64(1), "wallet": "0xabc",
				"points":                  "1000",
				"percent_share_of_points": "0.05",
				"total_volume":            "10000000",
			},
		},
	})
	got, err := api.GetPointsLeaderboard(context.Background(), types.PointsLeaderboardQuery{Program: "trading", Page: 1})
	require.NoError(t, err)
	require.Len(t, got.Leaderboard, 1)
	assert.Equal(t, int64(1), got.Leaderboard[0].Rank)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, float64(1), params["page"])
}
