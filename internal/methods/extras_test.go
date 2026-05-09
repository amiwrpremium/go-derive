package methods_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

// publicMethods covers every public/* extra wrapper. They all share the
// same shape: take params, return json.RawMessage. Param-shape correctness
// is enforced server-side; the SDK only forwards.
func TestExtras_PublicMethods(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	cases := []struct {
		name    string
		method  string
		invoke  func() (json.RawMessage, error)
		params  map[string]any
		mockOut any
	}{
		{
			name:   "GetFundingRateHistory",
			method: "public/get_funding_rate_history",
			invoke: func() (json.RawMessage, error) {
				return api.GetFundingRateHistory(context.Background(), map[string]any{"instrument_name": "BTC-PERP"})
			},
			mockOut: []map[string]any{{"timestamp": 1, "funding_rate": "0.0001"}},
		},
		{
			name:   "GetPerpImpactTWAP",
			method: "public/get_perp_impact_twap",
			invoke: func() (json.RawMessage, error) {
				return api.GetPerpImpactTWAP(context.Background(), map[string]any{
					"currency": "BTC", "start_time": 0, "end_time": 1,
				})
			},
			mockOut: map[string]any{"impact_price": "100"},
		},
		{
			name:   "GetLatestSignedFeeds",
			method: "public/get_latest_signed_feeds",
			invoke: func() (json.RawMessage, error) {
				return api.GetLatestSignedFeeds(context.Background(), nil)
			},
			mockOut: map[string]any{"feeds": []any{}},
		},
		{
			name:   "GetSpotFeedHistory",
			method: "public/get_spot_feed_history",
			invoke: func() (json.RawMessage, error) {
				return api.GetSpotFeedHistory(context.Background(), map[string]any{
					"currency": "BTC", "period": "1h",
					"start_timestamp": 0, "end_timestamp": 1,
				})
			},
			mockOut: []any{},
		},
		{
			name:   "GetStatistics",
			method: "public/statistics",
			invoke: func() (json.RawMessage, error) {
				return api.GetStatistics(context.Background(), "BTC-PERP")
			},
			mockOut: map[string]any{"open_interest": "1"},
		},
		{
			name:   "GetTransaction",
			method: "public/get_transaction",
			invoke: func() (json.RawMessage, error) {
				return api.GetTransaction(context.Background(), "tx-1")
			},
			mockOut: map[string]any{"status": "settled"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ft.HandleResult(c.method, c.mockOut)
			raw, err := c.invoke()
			require.NoError(t, err)
			assert.NotEmpty(t, raw)
		})
	}
}

// privateMethods covers the read + write extras that require a signer.
// Each test seeds the response, invokes, and asserts the route matched.
func TestExtras_PrivateMethods(t *testing.T) {
	api, ft := newAPI(t, true, 9)
	cases := []struct {
		name    string
		method  string
		invoke  func() (json.RawMessage, error)
		mockOut any
	}{
		{"Replace", "private/replace",
			func() (json.RawMessage, error) {
				return api.Replace(context.Background(),
					map[string]any{"order_id_to_cancel": "abc"})
			}, map[string]any{}},
		{"OrderDebug", "private/order_debug",
			func() (json.RawMessage, error) {
				return api.OrderDebug(context.Background(),
					map[string]any{"instrument_name": "BTC-PERP"})
			}, map[string]any{}},
		{"CancelByNonce", "private/cancel_by_nonce",
			func() (json.RawMessage, error) {
				return api.CancelByNonce(context.Background(), "BTC-PERP", 42)
			}, map[string]any{}},
		{"SetCancelOnDisconnect", "private/set_cancel_on_disconnect",
			func() (json.RawMessage, error) {
				return api.SetCancelOnDisconnect(context.Background(), true)
			}, map[string]any{}},
		{"ChangeSubaccountLabel", "private/change_subaccount_label",
			func() (json.RawMessage, error) {
				return api.ChangeSubaccountLabel(context.Background(), "alpha")
			}, map[string]any{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ft.HandleResult(c.method, c.mockOut)
			raw, err := c.invoke()
			require.NoError(t, err, "method %s", c.method)
			assert.NotEmpty(t, raw)
		})
	}
}

// Without a signer, every private extra returns ErrUnauthorized.
func TestExtras_PrivateRequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	checks := []func() (json.RawMessage, error){
		func() (json.RawMessage, error) { return api.OrderDebug(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.Replace(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.SetCancelOnDisconnect(context.Background(), true) },
	}
	for _, fn := range checks {
		_, err := fn()
		assert.ErrorIs(t, err, derrors.ErrUnauthorized)
	}
}
