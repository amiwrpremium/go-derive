package derive_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// newAPI returns a *derive.API wired to a FakeTransport.
// signed=true attaches a LocalSigner and a default subaccount id.
func newAPI(t *testing.T, signed bool, sub int64) (*derive.API, *testutil.FakeTransport) {
	t.Helper()
	ft := testutil.NewFakeTransport()
	api := &derive.API{
		T:               ft,
		Domain:          derive.Mainnet().EIP712Domain(),
		Nonces:          derive.NewNonceGen(),
		SignatureExpiry: 300,
	}
	api.SetTradeModule(common.HexToAddress(derive.Mainnet().Contracts.TradeModule))
	if signed {
		s, err := derive.NewLocalSigner(testKey)
		require.NoError(t, err)
		api.Signer = s
		api.Subaccount = sub
	}
	return api, ft
}

func paramsAsMap(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()
	if len(raw) == 0 {
		return nil
	}
	var m map[string]any
	require.NoError(t, json.Unmarshal(raw, &m))
	return m
}
func TestGetCollateral_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_collaterals", map[string]any{
		"collaterals": []map[string]any{
			{"asset_name": "USDC", "asset_type": "erc20", "amount": "100", "mark_value": "100"},
		},
	})
	got, err := api.GetCollateral(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "USDC", got[0].AssetName)
	assert.Equal(t, "private/get_collaterals", ft.LastCall().Method)
}

func TestGetCollateral_Empty(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_collaterals", map[string]any{"collaterals": []any{}})
	got, err := api.GetCollateral(context.Background())
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestGetCollateral_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.GetCollateral(context.Background())
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}

// boom is the error injected by the fake transport on every wrapper under
// test below. Using a single sentinel keeps the table compact.
var boom = &derive.APIError{Code: 9999, Message: "boom"}

// rawWrapper is a generic adapter for the family of wrappers that take
// (ctx, map[string]any) and return (json.RawMessage, error). It lets us
// table-drive every such wrapper through one test.
type rawWrapper func(ctx context.Context, params map[string]any) (json.RawMessage, error)

func TestRPCWrappers_PropagateError(t *testing.T) {
	cases := []struct {
		name   string
		method string
		needs  int64
		fn     func(*derive.API) rawWrapper
	}{

		{"GetFundingHistory", "private/get_funding_history", 1, func(a *derive.API) rawWrapper { return a.GetFundingHistory }},
		{"GetLiquidationHistory", "private/get_liquidation_history", 1, func(a *derive.API) rawWrapper { return a.GetLiquidationHistory }},
		{"GetOptionSettlementHistory", "private/get_option_settlement_history", 1, func(a *derive.API) rawWrapper { return a.GetOptionSettlementHistory }},
		{"GetSubaccountValueHistory", "private/get_subaccount_value_history", 1, func(a *derive.API) rawWrapper { return a.GetSubaccountValueHistory }},
		{"GetERC20TransferHistory", "private/get_erc20_transfer_history", 1, func(a *derive.API) rawWrapper { return a.GetERC20TransferHistory }},
		{"GetInterestHistory", "private/get_interest_history", 1, func(a *derive.API) rawWrapper { return a.GetInterestHistory }},
		{"ExpiredAndCancelledHistory", "private/expired_and_cancelled_history", 1, func(a *derive.API) rawWrapper { return a.ExpiredAndCancelledHistory }},
		{"GetNotifications", "private/get_notifications", 1, func(a *derive.API) rawWrapper { return a.GetNotifications }},
		{"UpdateNotifications", "private/update_notifications", 1, func(a *derive.API) rawWrapper { return a.UpdateNotifications }},
		{"Replace", "private/replace", 1, func(a *derive.API) rawWrapper { return a.Replace }},
		{"OrderDebug", "private/order_debug", 1, func(a *derive.API) rawWrapper { return a.OrderDebug }},

		{"GetRFQs", "private/get_rfqs", 1, func(a *derive.API) rawWrapper { return a.GetRFQs }},
		{"GetQuotes", "private/get_quotes", 1, func(a *derive.API) rawWrapper { return a.GetQuotes }},
		{"PollQuotes", "private/poll_quotes", 1, func(a *derive.API) rawWrapper { return a.PollQuotes }},
		{"SendQuote", "private/send_quote", 1, func(a *derive.API) rawWrapper { return a.SendQuote }},
		{"ExecuteQuote", "private/execute_quote", 1, func(a *derive.API) rawWrapper { return a.ExecuteQuote }},
		{"CancelQuote", "private/cancel_quote", 1, func(a *derive.API) rawWrapper { return a.CancelQuote }},
		{"CancelBatchQuotes", "private/cancel_batch_quotes", 1, func(a *derive.API) rawWrapper { return a.CancelBatchQuotes }},
		{"CancelBatchRFQs", "private/cancel_batch_rfqs", 1, func(a *derive.API) rawWrapper { return a.CancelBatchRFQs }},
		{"RFQGetBestQuote", "private/rfq_get_best_quote", 1, func(a *derive.API) rawWrapper { return a.RFQGetBestQuote }},
		{"OrderQuote", "private/order_quote", 1, func(a *derive.API) rawWrapper { return a.OrderQuote }},

		{"GetFundingRateHistory", "public/get_funding_rate_history", 0, func(a *derive.API) rawWrapper { return a.GetFundingRateHistory }},
		{"GetPerpImpactTWAP", "public/get_perp_impact_twap", 0, func(a *derive.API) rawWrapper { return a.GetPerpImpactTWAP }},
		{"GetPublicMargin", "public/get_margin", 0, func(a *derive.API) rawWrapper { return a.GetPublicMargin }},
		{"GetLatestSignedFeeds", "public/get_latest_signed_feeds", 0, func(a *derive.API) rawWrapper { return a.GetLatestSignedFeeds }},
		{"GetSpotFeedHistory", "public/get_spot_feed_history", 0, func(a *derive.API) rawWrapper { return a.GetSpotFeedHistory }},
		{"GetPublicOptionSettlementHistory", "public/get_option_settlement_history", 0, func(a *derive.API) rawWrapper { return a.GetPublicOptionSettlementHistory }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			api, ft := newAPI(t, true, c.needs)
			ft.HandleError(c.method, boom)
			_, err := c.fn(api)(context.Background(), map[string]any{})
			assert.Error(t, err)
			var apiErr *derive.APIError
			assert.True(t, errors.As(err, &apiErr))
			assert.Equal(t, 9999, apiErr.Code)
		})
	}
}

// TestNoArgWrappers_PropagateError covers wrappers that don't take a
// params map and so don't fit rawWrapper.
func TestNoArgWrappers_PropagateError(t *testing.T) {
	t.Run("GetMargin", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_margin", boom)
		_, err := api.GetMargin(context.Background())
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("GetMMPConfig", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_mmp_config", boom)
		_, err := api.GetMMPConfig(context.Background())
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("GetAccount", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_account", boom)
		_, err := api.GetAccount(context.Background())
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("CancelByNonce", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/cancel_by_nonce", boom)
		_, err := api.CancelByNonce(context.Background(), "BTC-PERP", 42)
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("SetCancelOnDisconnect", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/set_cancel_on_disconnect", boom)
		_, err := api.SetCancelOnDisconnect(context.Background(), true)
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("ChangeSubaccountLabel", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/change_subaccount_label", boom)
		_, err := api.ChangeSubaccountLabel(context.Background(), "newlabel")
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("GetStatistics", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/statistics", boom)
		_, err := api.GetStatistics(context.Background(), "BTC-PERP")
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("GetTransaction", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/get_transaction", boom)
		_, err := api.GetTransaction(context.Background(), "TX1")
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("GetCurrencies", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/get_all_currencies", boom)
		_, err := api.GetCurrencies(context.Background())
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("CancelByLabel", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/cancel_by_label", boom)
		_, err := api.CancelByLabel(context.Background(), "L")
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("CancelByInstrument", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/cancel_by_instrument", boom)
		_, err := api.CancelByInstrument(context.Background(), "BTC-PERP")
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("CancelAll", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/cancel_all", boom)
		_, err := api.CancelAll(context.Background())
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("GetDepositHistory_ServerError", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_deposit_history", boom)
		_, _, err := api.GetDepositHistory(context.Background(), derive.PageRequest{Page: 1, PageSize: 10})
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("GetWithdrawalHistory_ServerError", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_withdrawal_history", boom)
		_, _, err := api.GetWithdrawalHistory(context.Background(), derive.PageRequest{})
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
	t.Run("GetTradeHistory_ServerError", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_trade_history", boom)
		_, _, err := api.GetTradeHistory(context.Background(), derive.PageRequest{})
		assert.ErrorAs(t, err, new(*derive.APIError))
	})
}

// silence unused-import warning when none of the testutil symbols
// surface in this file (newAPI lives in the same package).
var _ = testutil.NewFakeTransport

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
			name:   "GetPublicMargin",
			method: "public/get_margin",
			invoke: func() (json.RawMessage, error) {
				return api.GetPublicMargin(context.Background(), map[string]any{
					"simulated_collaterals": []any{},
					"simulated_positions":   []any{},
					"margin_type":           "PM",
				})
			},
			mockOut: map[string]any{"initial_margin": "0", "maintenance_margin": "0"},
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
		{
			name:   "GetPublicOptionSettlementHistory",
			method: "public/get_option_settlement_history",
			invoke: func() (json.RawMessage, error) {
				return api.GetPublicOptionSettlementHistory(context.Background(), nil)
			},
			mockOut: []any{},
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
		{"GetAccount", "private/get_account",
			func() (json.RawMessage, error) { return api.GetAccount(context.Background()) },
			map[string]any{"wallet": "0x"}},
		{"GetMargin", "private/get_margin",
			func() (json.RawMessage, error) { return api.GetMargin(context.Background()) },
			map[string]any{"initial_margin": "0"}},
		{"GetFundingHistory", "private/get_funding_history",
			func() (json.RawMessage, error) { return api.GetFundingHistory(context.Background(), nil) },
			[]any{}},
		{"GetLiquidationHistory", "private/get_liquidation_history",
			func() (json.RawMessage, error) { return api.GetLiquidationHistory(context.Background(), nil) },
			[]any{}},
		{"GetOptionSettlementHistory", "private/get_option_settlement_history",
			func() (json.RawMessage, error) {
				return api.GetOptionSettlementHistory(context.Background(), nil)
			},
			[]any{}},
		{"GetSubaccountValueHistory", "private/get_subaccount_value_history",
			func() (json.RawMessage, error) {
				return api.GetSubaccountValueHistory(context.Background(),
					map[string]any{"period": "1h", "start_timestamp": 0, "end_timestamp": 1})
			},
			[]any{}},
		{"GetERC20TransferHistory", "private/get_erc20_transfer_history",
			func() (json.RawMessage, error) {
				return api.GetERC20TransferHistory(context.Background(), nil)
			}, []any{}},
		{"GetInterestHistory", "private/get_interest_history",
			func() (json.RawMessage, error) {
				return api.GetInterestHistory(context.Background(), nil)
			}, []any{}},
		{"ExpiredAndCancelledHistory", "private/expired_and_cancelled_history",
			func() (json.RawMessage, error) {
				return api.ExpiredAndCancelledHistory(context.Background(), nil)
			}, []any{}},
		{"GetMMPConfig", "private/get_mmp_config",
			func() (json.RawMessage, error) { return api.GetMMPConfig(context.Background()) },
			map[string]any{}},
		{"GetNotifications", "private/get_notifications",
			func() (json.RawMessage, error) {
				return api.GetNotifications(context.Background(), nil)
			}, []any{}},
		{"UpdateNotifications", "private/update_notifications",
			func() (json.RawMessage, error) {
				return api.UpdateNotifications(context.Background(),
					map[string]any{"notification_ids": []int{1}, "status": "seen"})
			}, map[string]any{}},
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
		func() (json.RawMessage, error) { return api.GetAccount(context.Background()) },
		func() (json.RawMessage, error) { return api.GetMargin(context.Background()) },
		func() (json.RawMessage, error) { return api.GetFundingHistory(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.GetLiquidationHistory(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.OrderDebug(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.Replace(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.SetCancelOnDisconnect(context.Background(), true) },
	}
	for _, fn := range checks {
		_, err := fn()
		assert.ErrorIs(t, err, derive.ErrUnauthorized)
	}
}
func TestGetInstruments(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_instruments", []map[string]any{
		{"instrument_name": "BTC-PERP", "instrument_type": "perp", "is_active": true,
			"base_currency": "BTC", "quote_currency": "USDC",
			"tick_size": "0.5", "minimum_amount": "0.001", "maximum_amount": "100", "amount_step": "0.001"},
	})

	insts, err := api.GetInstruments(context.Background(), "BTC", derive.InstrumentTypePerp)
	require.NoError(t, err)
	require.Len(t, insts, 1)
	assert.Equal(t, "BTC-PERP", insts[0].Name)

	last := ft.LastCall()
	assert.Equal(t, "public/get_instruments", last.Method)
	params := paramsAsMap(t, last.Params)
	assert.Equal(t, "BTC", params["currency"])
	assert.Equal(t, "perp", params["instrument_type"])
	assert.Equal(t, false, params["expired"])
}

func TestGetInstruments_NoFilters(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_instruments", []map[string]any{})

	_, err := api.GetInstruments(context.Background(), "", "")
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	_, hasCurrency := params["currency"]
	_, hasKind := params["instrument_type"]
	assert.False(t, hasCurrency)
	assert.False(t, hasKind)
}

func TestGetInstrument(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_instrument", map[string]any{
		"instrument_name": "BTC-PERP",
		"instrument_type": "perp",
		"is_active":       true,
		"base_currency":   "BTC",
		"quote_currency":  "USDC",
		"tick_size":       "0.5",
		"minimum_amount":  "0.001",
		"maximum_amount":  "100",
		"amount_step":     "0.001",
	})

	got, err := api.GetInstrument(context.Background(), "BTC-PERP")
	require.NoError(t, err)
	assert.Equal(t, "BTC-PERP", got.Name)
}

func TestGetTicker(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_ticker", map[string]any{
		"instrument_name": "BTC-PERP",
		"best_bid_price":  "100", "best_bid_amount": "1",
		"best_ask_price": "101", "best_ask_amount": "2",
		"mark_price":  "100.5",
		"index_price": "100.5",
		"timestamp":   1700000000000,
	})
	got, err := api.GetTicker(context.Background(), "BTC-PERP")
	require.NoError(t, err)
	assert.Equal(t, "BTC-PERP", got.InstrumentName)
}

func TestGetPublicTradeHistory(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_trade_history", map[string]any{
		"trades":     []any{},
		"pagination": map[string]any{"num_pages": 1, "count": 0, "current_page": 1, "page_size": 50},
	})
	trades, page, err := api.GetPublicTradeHistory(context.Background(), "BTC-PERP", derive.PageRequest{Page: 2, PageSize: 50})
	require.NoError(t, err)
	assert.Empty(t, trades)
	assert.Equal(t, 1, page.NumPages)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, float64(2), params["page"])
	assert.Equal(t, float64(50), params["page_size"])
}

func TestGetTime(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_time", 1700000000000)
	got, err := api.GetTime(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(1700000000000), got)
}

func TestGetCurrencies(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_all_currencies", []map[string]any{
		{"currency": "USDC"}, {"currency": "WETH"},
	})
	got, err := api.GetCurrencies(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"USDC", "WETH"}, got)
}

func TestPublicMethods_PropagateUnhandledTransportError(t *testing.T) {
	api, _ := newAPI(t, false, 0)

	_, err := api.GetTime(context.Background())
	require.Error(t, err)
}

// silence unused json import warning on platforms where compiler is strict.
var _ = json.Marshal

func TestMMPConfig_Validate_Happy(t *testing.T) {
	cfg := derive.MMPConfig{Currency: "BTC", MMPFrozenTimeMs: 1000, MMPIntervalMs: 500}
	require.NoError(t, cfg.Validate())
}

func TestMMPConfig_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		cfg  derive.MMPConfig
		want string
	}{
		{"empty currency", derive.MMPConfig{}, "currency"},
		{"negative frozen", derive.MMPConfig{Currency: "BTC", MMPFrozenTimeMs: -1}, "mmp_frozen_time"},
		{"negative interval", derive.MMPConfig{Currency: "BTC", MMPIntervalMs: -1}, "mmp_interval"},
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
	cfg := derive.MMPConfig{
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
	cfg := derive.MMPConfig{Currency: "BTC", MMPDeltaLimit: "1"}
	require.NoError(t, api.SetMMPConfig(context.Background(), cfg))
	params := paramsAsMap(t, ft.LastCall().Params)
	_, has := params["mmp_amount_limit"]
	assert.False(t, has)
}

func TestSetMMPConfig_OmitsEmptyDeltaLimit(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/set_mmp_config", nil)
	cfg := derive.MMPConfig{Currency: "BTC", MMPAmountLimit: "1"}
	require.NoError(t, api.SetMMPConfig(context.Background(), cfg))
	params := paramsAsMap(t, ft.LastCall().Params)
	_, has := params["mmp_delta_limit"]
	assert.False(t, has)
}

func TestSetMMPConfig_OmitsBothLimits(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/set_mmp_config", nil)
	cfg := derive.MMPConfig{Currency: "BTC"}
	require.NoError(t, api.SetMMPConfig(context.Background(), cfg))
	params := paramsAsMap(t, ft.LastCall().Params)
	_, hasA := params["mmp_amount_limit"]
	_, hasD := params["mmp_delta_limit"]
	assert.False(t, hasA)
	assert.False(t, hasD)
}

func TestSetMMPConfig_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	err := api.SetMMPConfig(context.Background(), derive.MMPConfig{Currency: "BTC"})
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
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
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}
func validPlaceOrderInput() derive.PlaceOrderInput {
	return derive.PlaceOrderInput{
		InstrumentName: "BTC-PERP",
		Asset:          common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Direction:      derive.DirectionBuy,
		OrderType:      derive.OrderTypeLimit,
		Amount:         derive.MustDecimal("1"),
		LimitPrice:     derive.MustDecimal("100"),
		MaxFee:         derive.MustDecimal("1"),
	}
}

func TestPlaceOrderInput_Validate_Happy(t *testing.T) {
	require.NoError(t, validPlaceOrderInput().Validate())
}

func TestPlaceOrderInput_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*derive.PlaceOrderInput)
		want string
	}{
		{"empty instrument", func(in *derive.PlaceOrderInput) { in.InstrumentName = "" }, "instrument_name"},
		{"zero asset", func(in *derive.PlaceOrderInput) { in.Asset = common.Address{} }, "asset"},
		{"bad direction", func(in *derive.PlaceOrderInput) { in.Direction = derive.Direction("x") }, "direction"},
		{"bad order type", func(in *derive.PlaceOrderInput) { in.OrderType = derive.OrderType("x") }, "order_type"},
		{"bad time-in-force", func(in *derive.PlaceOrderInput) { in.TimeInForce = derive.TimeInForce("x") }, "time_in_force"},
		{"zero amount", func(in *derive.PlaceOrderInput) { in.Amount = derive.MustDecimal("0") }, "amount"},
		{"zero price", func(in *derive.PlaceOrderInput) { in.LimitPrice = derive.MustDecimal("0") }, "limit_price"},
		{"negative fee", func(in *derive.PlaceOrderInput) { in.MaxFee = derive.MustDecimal("-1") }, "max_fee"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in := validPlaceOrderInput()
			c.mut(&in)
			err := in.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, derive.ErrInvalidParams))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func TestPlaceOrderInput_Validate_AllowsEmptyTimeInForce(t *testing.T) {
	in := validPlaceOrderInput()
	in.TimeInForce = ""
	require.NoError(t, in.Validate())
}

func TestPlaceOrder_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.PlaceOrder(context.Background(), derive.PlaceOrderInput{})
	assert.ErrorIs(t, err, derive.ErrUnauthorized)
}

func TestPlaceOrder_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.PlaceOrder(context.Background(), derive.PlaceOrderInput{})
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}

func TestPlaceOrder_Success_PopulatesSignatureFields(t *testing.T) {
	api, ft := newAPI(t, true, 1)

	ft.Handle("private/order", func(_ json.RawMessage) (any, error) {

		return map[string]any{
			"order": map[string]any{
				"order_id": "O1", "subaccount_id": 1, "instrument_name": "BTC-PERP",
				"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
				"order_status": "open", "amount": "0.1", "filled_amount": "0",
				"limit_price": "65000", "max_fee": "10", "nonce": 1,
				"signer":             "0x0000000000000000000000000000000000000000",
				"creation_timestamp": 1700000000000, "last_update_timestamp": 1700000000000,
			},
		}, nil
	})

	in := derive.PlaceOrderInput{
		InstrumentName: "BTC-PERP",
		Asset:          common.HexToAddress("0x1111111111111111111111111111111111111111"),
		SubID:          0,
		Direction:      derive.DirectionBuy,
		OrderType:      derive.OrderTypeLimit,
		TimeInForce:    derive.TimeInForceGTC,
		Amount:         derive.MustDecimal("0.1"),
		LimitPrice:     derive.MustDecimal("65000"),
		MaxFee:         derive.MustDecimal("10"),
	}
	order, err := api.PlaceOrder(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, "O1", order.OrderID)

	params := paramsAsMap(t, ft.LastCall().Params)
	sig, _ := params["signature"].(string)
	assert.True(t, strings.HasPrefix(sig, "0x") && len(sig) == 132, "signature: %s", sig)
	assert.NotEmpty(t, params["signer"])
	assert.Greater(t, params["signature_expiry_sec"], float64(0))
	assert.Greater(t, params["nonce"], float64(0))
	assert.Equal(t, "buy", params["direction"])
	assert.Equal(t, "BTC-PERP", params["instrument_name"])
}

func TestCancelOrder(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/cancel", nil)
	require.NoError(t, api.CancelOrder(context.Background(), "BTC-PERP", "O1"))
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "BTC-PERP", params["instrument_name"])
	assert.Equal(t, "O1", params["order_id"])
	assert.Equal(t, float64(1), params["subaccount_id"])
}

func TestCancelOrder_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	err := api.CancelOrder(context.Background(), "BTC-PERP", "O1")
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}

func TestCancelByLabel(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/cancel_by_label", map[string]any{"cancelled_orders": 3})
	n, err := api.CancelByLabel(context.Background(), "alpha")
	require.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, "private/cancel_by_label", ft.LastCall().Method)
}

func TestCancelByInstrument(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/cancel_by_instrument", map[string]any{"cancelled_orders": 2})
	n, err := api.CancelByInstrument(context.Background(), "BTC-PERP")
	require.NoError(t, err)
	assert.Equal(t, 2, n)
}

func TestCancelAll(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/cancel_all", map[string]any{"cancelled_orders": 5})
	n, err := api.CancelAll(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 5, n)
}

func TestGetOrder(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_order", map[string]any{
		"order": map[string]any{
			"order_id": "O1", "subaccount_id": 1, "instrument_name": "BTC-PERP",
			"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
			"order_status": "open", "amount": "0.1", "filled_amount": "0",
			"limit_price": "65000", "max_fee": "10", "nonce": 1,
			"signer":             "0x0000000000000000000000000000000000000000",
			"creation_timestamp": 1, "last_update_timestamp": 1,
		},
	})
	o, err := api.GetOrder(context.Background(), "O1")
	require.NoError(t, err)
	assert.Equal(t, "O1", o.OrderID)
}

func TestGetOpenOrders(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_open_orders", map[string]any{"orders": []any{}})
	got, err := api.GetOpenOrders(context.Background())
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestGetOrderHistory(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_orders", map[string]any{
		"orders":     []any{},
		"pagination": map[string]any{"num_pages": 2, "count": 100},
	})
	_, page, err := api.GetOrderHistory(context.Background(), derive.PageRequest{Page: 1, PageSize: 50})
	require.NoError(t, err)
	assert.Equal(t, 2, page.NumPages)
	assert.Equal(t, 100, page.Count)
}

func TestPrivateMethods_RequireSubaccount_Across(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	cases := map[string]func() error{
		"GetOrder":      func() error { _, e := api.GetOrder(context.Background(), "x"); return e },
		"GetOpenOrders": func() error { _, e := api.GetOpenOrders(context.Background()); return e },
		"GetOrderHistory": func() error {
			_, _, e := api.GetOrderHistory(context.Background(), derive.PageRequest{})
			return e
		},
		"CancelByLabel":      func() error { _, e := api.CancelByLabel(context.Background(), "x"); return e },
		"CancelByInstrument": func() error { _, e := api.CancelByInstrument(context.Background(), "x"); return e },
		"CancelAll":          func() error { _, e := api.CancelAll(context.Background()); return e },
	}
	for name, fn := range cases {
		t.Run(name, func(t *testing.T) {
			assert.ErrorIs(t, fn(), derive.ErrSubaccountRequired)
		})
	}
}
func TestGetPositions_Empty(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_positions", map[string]any{"positions": []any{}})
	got, err := api.GetPositions(context.Background())
	require.NoError(t, err)
	assert.Empty(t, got)
	assert.Equal(t, "private/get_positions", ft.LastCall().Method)
}

func TestGetPositions_NonEmpty(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_positions", map[string]any{
		"positions": []map[string]any{{
			"instrument_name": "BTC-PERP",
			"instrument_type": "perp",
			"amount":          "0.5",
			"average_price":   "65000",
			"mark_price":      "65500",
			"mark_value":      "32750",
			"unrealized_pnl":  "250",
			"realized_pnl":    "0",
		}},
	})
	got, err := api.GetPositions(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "BTC-PERP", got[0].InstrumentName)
}

func TestGetPositions_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.GetPositions(context.Background())
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}
func TestSendRFQ_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/send_rfq", map[string]any{
		"rfq_id": "R1", "subaccount_id": 1, "status": "open",
		"legs": []any{}, "creation_timestamp": 1, "last_update_timestamp": 1,
	})
	rfq, err := api.SendRFQ(context.Background(), nil, derive.MustDecimal("100"))
	require.NoError(t, err)
	assert.Equal(t, "R1", rfq.RFQID)
}

func TestSendRFQ_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.SendRFQ(context.Background(), nil, derive.MustDecimal("0"))
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}

func TestPollRFQs_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/poll_rfqs", map[string]any{"rfqs": []any{}})
	got, err := api.PollRFQs(context.Background())
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestPollRFQs_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.PollRFQs(context.Background())
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}

func TestCancelRFQ_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/cancel_rfq", nil)
	require.NoError(t, api.CancelRFQ(context.Background(), "R1"))
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "R1", params["rfq_id"])
}

func TestCancelRFQ_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	err := api.CancelRFQ(context.Background(), "R1")
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}
func TestRFQExtras_AllMethods(t *testing.T) {
	api, ft := newAPI(t, true, 9)
	cases := []struct {
		name    string
		method  string
		invoke  func() (json.RawMessage, error)
		mockOut any
	}{
		{"GetRFQs", "private/get_rfqs",
			func() (json.RawMessage, error) { return api.GetRFQs(context.Background(), nil) },
			[]any{}},
		{"GetQuotes", "private/get_quotes",
			func() (json.RawMessage, error) { return api.GetQuotes(context.Background(), nil) },
			[]any{}},
		{"PollQuotes", "private/poll_quotes",
			func() (json.RawMessage, error) { return api.PollQuotes(context.Background(), nil) },
			[]any{}},
		{"SendQuote", "private/send_quote",
			func() (json.RawMessage, error) {
				return api.SendQuote(context.Background(), map[string]any{"rfq_id": "R"})
			},
			map[string]any{"quote_id": "Q"}},
		{"ExecuteQuote", "private/execute_quote",
			func() (json.RawMessage, error) {
				return api.ExecuteQuote(context.Background(), map[string]any{"quote_id": "Q"})
			},
			map[string]any{"status": "filled"}},
		{"CancelQuote", "private/cancel_quote",
			func() (json.RawMessage, error) {
				return api.CancelQuote(context.Background(), map[string]any{"quote_id": "Q"})
			},
			map[string]any{}},
		{"CancelBatchQuotes", "private/cancel_batch_quotes",
			func() (json.RawMessage, error) {
				return api.CancelBatchQuotes(context.Background(), nil)
			},
			map[string]any{}},
		{"CancelBatchRFQs", "private/cancel_batch_rfqs",
			func() (json.RawMessage, error) {
				return api.CancelBatchRFQs(context.Background(), nil)
			},
			map[string]any{}},
		{"RFQGetBestQuote", "private/rfq_get_best_quote",
			func() (json.RawMessage, error) {
				return api.RFQGetBestQuote(context.Background(), map[string]any{"rfq_id": "R"})
			},
			map[string]any{"price": "1"}},
		{"OrderQuote", "private/order_quote",
			func() (json.RawMessage, error) {
				return api.OrderQuote(context.Background(), map[string]any{"instrument_name": "BTC-PERP"})
			},
			map[string]any{}},
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

func TestRFQExtras_RequireSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	checks := []func() (json.RawMessage, error){
		func() (json.RawMessage, error) { return api.GetRFQs(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.GetQuotes(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.SendQuote(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.ExecuteQuote(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.OrderQuote(context.Background(), nil) },
	}
	for _, fn := range checks {
		_, err := fn()
		assert.ErrorIs(t, err, derive.ErrUnauthorized)
	}
}
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
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
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
	assert.ErrorIs(t, err, derive.ErrUnauthorized)
}
func TestGetTradeHistory_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_trade_history", map[string]any{
		"trades":     []any{},
		"pagination": map[string]any{"num_pages": 4, "count": 100},
	})
	_, page, err := api.GetTradeHistory(context.Background(), derive.PageRequest{Page: 1, PageSize: 25})
	require.NoError(t, err)
	assert.Equal(t, 4, page.NumPages)
	assert.Equal(t, 100, page.Count)

	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, float64(1), params["page"])
	assert.Equal(t, float64(25), params["page_size"])
}

func TestGetTradeHistory_OmitsZeroPagination(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_trade_history", map[string]any{
		"trades": []any{}, "pagination": map[string]any{},
	})
	_, _, err := api.GetTradeHistory(context.Background(), derive.PageRequest{})
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	_, hasPage := params["page"]
	_, hasSize := params["page_size"]
	assert.False(t, hasPage)
	assert.False(t, hasSize)
}

func TestGetTradeHistory_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, _, err := api.GetTradeHistory(context.Background(), derive.PageRequest{})
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}

func TestGetDepositHistory_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_deposit_history", map[string]any{
		"events":     []any{},
		"pagination": map[string]any{"num_pages": 0, "count": 0, "current_page": 1, "page_size": 10},
	})
	_, _, err := api.GetDepositHistory(context.Background(), derive.PageRequest{Page: 0, PageSize: 0})
	require.NoError(t, err)
}

func TestGetDepositHistory_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, _, err := api.GetDepositHistory(context.Background(), derive.PageRequest{})
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}

func TestGetWithdrawalHistory_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_withdrawal_history", map[string]any{
		"events":     []any{},
		"pagination": map[string]any{"num_pages": 0, "count": 0, "current_page": 1, "page_size": 10},
	})
	_, _, err := api.GetWithdrawalHistory(context.Background(), derive.PageRequest{Page: 2, PageSize: 50})
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, float64(2), params["page"])
	assert.Equal(t, float64(50), params["page_size"])
}

func TestGetWithdrawalHistory_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, _, err := api.GetWithdrawalHistory(context.Background(), derive.PageRequest{})
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}
