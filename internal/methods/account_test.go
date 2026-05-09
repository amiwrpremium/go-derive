package methods_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestGetAccount_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_account", map[string]any{
		"subaccount_ids":             []any{int64(1), int64(2)},
		"wallet":                     "0xabc",
		"cancel_on_disconnect":       true,
		"creation_timestamp_sec":     int64(1700000000),
		"is_rfq_maker":               false,
		"referral_code":              "REF",
		"websocket_matching_tps":     int64(50),
		"websocket_non_matching_tps": int64(20),
		"websocket_option_tps":       int64(10),
		"websocket_perp_tps":         int64(15),
		"per_endpoint_tps":           map[string]any{},
		"fee_info": map[string]any{
			"base_fee_discount":  "0",
			"option_maker_fee":   "0.0003",
			"option_taker_fee":   "0.0005",
			"perp_maker_fee":     "0",
			"perp_taker_fee":     "0.0003",
			"rfq_maker_discount": "0",
			"rfq_taker_discount": "0",
			"spot_maker_fee":     "0",
			"spot_taker_fee":     "0",
		},
	})
	got, err := api.GetAccount(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "0xabc", got.Wallet)
	assert.Equal(t, []int64{1, 2}, got.SubaccountIDs)
	assert.True(t, got.CancelOnDisconnect)
	assert.Equal(t, "0.0003", got.FeeInfo.OptionMakerFee.String())
}

func TestGetAccount_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.GetAccount(context.Background())
	assert.True(t, errors.Is(err, derrors.ErrUnauthorized))
}

func TestGetAccount_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleError("private/get_account", boom)
	_, err := api.GetAccount(context.Background())
	assert.ErrorAs(t, err, new(*derrors.APIError))
}

func TestGetMargin_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/get_margin", map[string]any{
		"subaccount_id":           int64(7),
		"is_valid_trade":          true,
		"pre_initial_margin":      "100",
		"post_initial_margin":     "110",
		"pre_maintenance_margin":  "50",
		"post_maintenance_margin": "55",
	})
	got, err := api.GetMargin(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(7), got.SubaccountID)
	assert.True(t, got.IsValidTrade)
	assert.Equal(t, "110", got.PostInitialMargin.String())
}

func TestGetMargin_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.GetMargin(context.Background())
	assert.True(t, errors.Is(err, derrors.ErrUnauthorized))
}

func TestGetMargin_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.GetMargin(context.Background())
	assert.True(t, errors.Is(err, derrors.ErrSubaccountRequired))
}

func TestGetMargin_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleError("private/get_margin", boom)
	_, err := api.GetMargin(context.Background())
	assert.ErrorAs(t, err, new(*derrors.APIError))
}

func TestGetPublicMargin_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_margin", map[string]any{
		"subaccount_id":           int64(0),
		"is_valid_trade":          false,
		"pre_initial_margin":      "0",
		"post_initial_margin":     "0",
		"pre_maintenance_margin":  "0",
		"post_maintenance_margin": "0",
	})
	got, err := api.GetPublicMargin(context.Background(), map[string]any{
		"margin_type":           "PM",
		"market":                "BTC",
		"simulated_collaterals": []any{},
		"simulated_positions":   []any{},
	})
	require.NoError(t, err)
	assert.False(t, got.IsValidTrade)
}

func TestGetPublicMargin_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleError("public/get_margin", boom)
	_, err := api.GetPublicMargin(context.Background(), map[string]any{})
	assert.ErrorAs(t, err, new(*derrors.APIError))
}
