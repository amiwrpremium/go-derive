package methods_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/methods"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestMMPConfig_Validate_Happy(t *testing.T) {
	cfg := methods.MMPConfig{Currency: "BTC", MMPFrozenTimeMs: 1000, MMPIntervalMs: 500}
	require.NoError(t, cfg.Validate())
}

func TestMMPConfig_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		cfg  methods.MMPConfig
		want string
	}{
		{"empty currency", methods.MMPConfig{}, "currency"},
		{"negative frozen", methods.MMPConfig{Currency: "BTC", MMPFrozenTimeMs: -1}, "mmp_frozen_time"},
		{"negative interval", methods.MMPConfig{Currency: "BTC", MMPIntervalMs: -1}, "mmp_interval"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.cfg.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, derive.ErrInvalidParams))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func TestSetMMPConfig_AllFieldsPopulated(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/set_mmp_config", nil)
	cfg := methods.MMPConfig{
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
	cfg := methods.MMPConfig{Currency: "BTC", MMPDeltaLimit: "1"}
	require.NoError(t, api.SetMMPConfig(context.Background(), cfg))
	params := paramsAsMap(t, ft.LastCall().Params)
	_, has := params["mmp_amount_limit"]
	assert.False(t, has)
}

func TestSetMMPConfig_OmitsEmptyDeltaLimit(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/set_mmp_config", nil)
	cfg := methods.MMPConfig{Currency: "BTC", MMPAmountLimit: "1"}
	require.NoError(t, api.SetMMPConfig(context.Background(), cfg))
	params := paramsAsMap(t, ft.LastCall().Params)
	_, has := params["mmp_delta_limit"]
	assert.False(t, has)
}

func TestSetMMPConfig_OmitsBothLimits(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/set_mmp_config", nil)
	cfg := methods.MMPConfig{Currency: "BTC"}
	require.NoError(t, api.SetMMPConfig(context.Background(), cfg))
	params := paramsAsMap(t, ft.LastCall().Params)
	_, hasA := params["mmp_amount_limit"]
	_, hasD := params["mmp_delta_limit"]
	assert.False(t, hasA)
	assert.False(t, hasD)
}

func TestSetMMPConfig_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	err := api.SetMMPConfig(context.Background(), methods.MMPConfig{Currency: "BTC"})
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
