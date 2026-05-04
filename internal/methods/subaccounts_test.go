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
		"subaccount_id":        7,
		"owner_address":        "0x1111111111111111111111111111111111111111",
		"margin_type":          "PM",
		"is_under_liquidation": false,
		"subaccount_value":     "0",
		"initial_margin":       "0",
		"maintenance_margin":   "0",
	})
	got, err := api.GetSubaccount(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(7), got.SubaccountID)

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
