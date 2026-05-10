package methods_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestMMPConfig_Validate_Happy(t *testing.T) {
	cfg := types.MMPConfig{Currency: "BTC", MMPFrozenTimeMs: 1000, MMPIntervalMs: 500}
	require.NoError(t, cfg.Validate())
}

func TestMMPConfig_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		cfg  types.MMPConfig
		want string
	}{
		{"empty currency", types.MMPConfig{}, "currency"},
		{"negative frozen", types.MMPConfig{Currency: "BTC", MMPFrozenTimeMs: -1}, "mmp_frozen_time"},
		{"negative interval", types.MMPConfig{Currency: "BTC", MMPIntervalMs: -1}, "mmp_interval"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.cfg.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, types.ErrInvalidParams))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func TestSetMMPConfig_AllFieldsPopulated(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/set_mmp_config", nil)
	cfg := types.MMPConfig{
		Currency:        "BTC",
		MMPFrozenTimeMs: 1000,
		MMPIntervalMs:   500,
		MMPAmountLimit:  "1",
		MMPDeltaLimit:   "0.5",
	}
	require.NoError(t, api.SetMMPConfig(context.Background(), cfg))
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "BTC", params["currency"])
	assert.Equal(t, "1", params["mmp_amount_limit"])
	assert.Equal(t, "0.5", params["mmp_delta_limit"])
	assert.Equal(t, float64(1000), params["mmp_frozen_time"])
	assert.Equal(t, float64(500), params["mmp_interval"])
}

func TestSetMMPConfig_OmitsEmptyAmountLimit(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/set_mmp_config", nil)
	cfg := types.MMPConfig{Currency: "BTC", MMPDeltaLimit: "1"}
	require.NoError(t, api.SetMMPConfig(context.Background(), cfg))
	params := paramsAsMap(t, ft.LastCall().Params)
	_, has := params["mmp_amount_limit"]
	assert.False(t, has)
}

func TestSetMMPConfig_OmitsEmptyDeltaLimit(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/set_mmp_config", nil)
	cfg := types.MMPConfig{Currency: "BTC", MMPAmountLimit: "1"}
	require.NoError(t, api.SetMMPConfig(context.Background(), cfg))
	params := paramsAsMap(t, ft.LastCall().Params)
	_, has := params["mmp_delta_limit"]
	assert.False(t, has)
}

func TestSetMMPConfig_OmitsBothLimits(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/set_mmp_config", nil)
	cfg := types.MMPConfig{Currency: "BTC"}
	require.NoError(t, api.SetMMPConfig(context.Background(), cfg))
	params := paramsAsMap(t, ft.LastCall().Params)
	_, hasA := params["mmp_amount_limit"]
	_, hasD := params["mmp_delta_limit"]
	assert.False(t, hasA)
	assert.False(t, hasD)
}

func TestSetMMPConfig_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	err := api.SetMMPConfig(context.Background(), types.MMPConfig{Currency: "BTC"})
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestResetMMP_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/reset_mmp", nil)
	require.NoError(t, api.ResetMMP(context.Background(), "BTC"))
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "BTC", params["currency"])
	assert.Equal(t, float64(1), params["subaccount_id"])
}

func TestResetMMP_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	err := api.ResetMMP(context.Background(), "BTC")
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestGetMMPConfig_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/get_mmp_config", []any{
		map[string]any{
			"subaccount_id":     int64(7),
			"currency":          "BTC",
			"mmp_frozen_time":   int64(5000),
			"mmp_interval":      int64(1000),
			"mmp_amount_limit":  "100",
			"mmp_delta_limit":   "50",
			"mmp_unfreeze_time": int64(0),
			"is_frozen":         false,
		},
	})
	got, err := api.GetMMPConfig(context.Background(), "")
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "BTC", got[0].Currency)
	assert.Equal(t, "100", got[0].MMPAmountLimit.String())
	assert.False(t, got[0].IsFrozen)
}

func TestGetMMPConfig_FilterByCurrency(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/get_mmp_config", []any{})
	_, err := api.GetMMPConfig(context.Background(), "ETH")
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "ETH", params["currency"])
}

func TestGetMMPConfig_OmitsEmptyCurrency(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/get_mmp_config", []any{})
	_, err := api.GetMMPConfig(context.Background(), "")
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	_, has := params["currency"]
	assert.False(t, has, "empty currency must be omitted from the params map")
}

func TestGetMMPConfig_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.GetMMPConfig(context.Background(), "")
	assert.True(t, errors.Is(err, derrors.ErrSubaccountRequired))
}

func TestGetMMPConfig_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleError("private/get_mmp_config", boom)
	_, err := api.GetMMPConfig(context.Background(), "")
	assert.ErrorAs(t, err, new(*derrors.APIError))
}
