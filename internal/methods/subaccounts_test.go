package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestGetSubaccount_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/get_subaccount", map[string]any{
		"subaccount_id":                  7,
		"currency":                       "USDC",
		"label":                          "main",
		"margin_type":                    "PM",
		"is_under_liquidation":           false,
		"subaccount_value":               "0",
		"initial_margin":                 "0",
		"maintenance_margin":             "0",
		"open_orders_margin":             "0",
		"projected_margin_change":        "0",
		"collaterals_initial_margin":     "0",
		"collaterals_maintenance_margin": "0",
		"collaterals_value":              "0",
		"positions_initial_margin":       "0",
		"positions_maintenance_margin":   "0",
		"positions_value":                "0",
		"collaterals":                    []any{},
		"open_orders":                    []any{},
		"positions":                      []any{},
	})
	got, err := api.GetSubaccount(context.Background())
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(7), got.SubaccountID)
	assert.Equal(t, "main", got.Label)
	assert.Equal(t, "USDC", got.Currency)

	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, float64(7), params["subaccount_id"])
}

func TestGetSubaccount_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.GetSubaccount(context.Background())
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestGetSubaccounts_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_subaccounts", map[string]any{"subaccount_ids": []int64{1, 2, 3}})
	got, err := api.GetSubaccounts(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []int64{1, 2, 3}, got)

	params := paramsAsMap(t, ft.LastCall().Params)
	assert.NotEmpty(t, params["wallet"])
}

func TestGetSubaccounts_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.GetSubaccounts(context.Background())
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
}

func TestChangeSubaccountLabel_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/change_subaccount_label", "ok")
	require.NoError(t, api.ChangeSubaccountLabel(context.Background(), "alpha"))
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "alpha", params["label"])
	assert.Equal(t, float64(7), params["subaccount_id"])
}

func TestChangeSubaccountLabel_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	err := api.ChangeSubaccountLabel(context.Background(), "x")
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
}

func TestChangeSubaccountLabel_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	err := api.ChangeSubaccountLabel(context.Background(), "x")
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}
