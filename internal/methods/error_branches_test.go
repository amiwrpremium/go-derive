package methods_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/internal/methods"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// boom is the error injected by the fake transport on every wrapper under
// test below. Using a single sentinel keeps the table compact.
var boom = &derrors.APIError{Code: 9999, Message: "boom"}

// rawWrapper is a generic adapter for the family of wrappers that take
// (ctx, map[string]any) and return (json.RawMessage, error). It lets us
// table-drive every such wrapper through one test.
type rawWrapper func(ctx context.Context, params map[string]any) (json.RawMessage, error)

func TestRPCWrappers_PropagateError(t *testing.T) {
	cases := []struct {
		name   string
		method string
		needs  int64
		fn     func(*methods.API) rawWrapper
	}{
		// extras.go — private map-of-any wrappers
		{"GetFundingHistory", "private/get_funding_history", 1, func(a *methods.API) rawWrapper { return a.GetFundingHistory }},
		{"GetLiquidationHistory", "private/get_liquidation_history", 1, func(a *methods.API) rawWrapper { return a.GetLiquidationHistory }},
		{"GetOptionSettlementHistory", "private/get_option_settlement_history", 1, func(a *methods.API) rawWrapper { return a.GetOptionSettlementHistory }},
		{"GetSubaccountValueHistory", "private/get_subaccount_value_history", 1, func(a *methods.API) rawWrapper { return a.GetSubaccountValueHistory }},
		{"GetERC20TransferHistory", "private/get_erc20_transfer_history", 1, func(a *methods.API) rawWrapper { return a.GetERC20TransferHistory }},
		{"GetInterestHistory", "private/get_interest_history", 1, func(a *methods.API) rawWrapper { return a.GetInterestHistory }},
		{"ExpiredAndCancelledHistory", "private/expired_and_cancelled_history", 1, func(a *methods.API) rawWrapper { return a.ExpiredAndCancelledHistory }},
		{"GetNotifications", "private/get_notifications", 1, func(a *methods.API) rawWrapper { return a.GetNotifications }},
		{"UpdateNotifications", "private/update_notifications", 1, func(a *methods.API) rawWrapper { return a.UpdateNotifications }},
		{"Replace", "private/replace", 1, func(a *methods.API) rawWrapper { return a.Replace }},
		{"OrderDebug", "private/order_debug", 1, func(a *methods.API) rawWrapper { return a.OrderDebug }},
		// rfq_extras.go
		{"GetRFQs", "private/get_rfqs", 1, func(a *methods.API) rawWrapper { return a.GetRFQs }},
		{"GetQuotes", "private/get_quotes", 1, func(a *methods.API) rawWrapper { return a.GetQuotes }},
		{"PollQuotes", "private/poll_quotes", 1, func(a *methods.API) rawWrapper { return a.PollQuotes }},
		{"SendQuote", "private/send_quote", 1, func(a *methods.API) rawWrapper { return a.SendQuote }},
		{"ExecuteQuote", "private/execute_quote", 1, func(a *methods.API) rawWrapper { return a.ExecuteQuote }},
		{"CancelQuote", "private/cancel_quote", 1, func(a *methods.API) rawWrapper { return a.CancelQuote }},
		{"CancelBatchQuotes", "private/cancel_batch_quotes", 1, func(a *methods.API) rawWrapper { return a.CancelBatchQuotes }},
		{"CancelBatchRFQs", "private/cancel_batch_rfqs", 1, func(a *methods.API) rawWrapper { return a.CancelBatchRFQs }},
		{"RFQGetBestQuote", "private/rfq_get_best_quote", 1, func(a *methods.API) rawWrapper { return a.RFQGetBestQuote }},
		{"OrderQuote", "private/order_quote", 1, func(a *methods.API) rawWrapper { return a.OrderQuote }},
		// public — map-of-any wrappers (no subaccount required)
		{"GetFundingRateHistory", "public/get_funding_rate_history", 0, func(a *methods.API) rawWrapper { return a.GetFundingRateHistory }},
		{"GetPerpImpactTWAP", "public/get_perp_impact_twap", 0, func(a *methods.API) rawWrapper { return a.GetPerpImpactTWAP }},
		{"GetPublicMargin", "public/get_margin", 0, func(a *methods.API) rawWrapper { return a.GetPublicMargin }},
		{"GetLatestSignedFeeds", "public/get_latest_signed_feeds", 0, func(a *methods.API) rawWrapper { return a.GetLatestSignedFeeds }},
		{"GetSpotFeedHistory", "public/get_spot_feed_history", 0, func(a *methods.API) rawWrapper { return a.GetSpotFeedHistory }},
		{"GetPublicOptionSettlementHistory", "public/get_option_settlement_history", 0, func(a *methods.API) rawWrapper { return a.GetPublicOptionSettlementHistory }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			api, ft := newAPI(t, true, c.needs)
			ft.HandleError(c.method, boom)
			_, err := c.fn(api)(context.Background(), map[string]any{})
			assert.Error(t, err)
			var apiErr *derrors.APIError
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
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetMMPConfig", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_mmp_config", boom)
		_, err := api.GetMMPConfig(context.Background())
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetAccount", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_account", boom)
		_, err := api.GetAccount(context.Background())
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("CancelByNonce", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/cancel_by_nonce", boom)
		_, err := api.CancelByNonce(context.Background(), "BTC-PERP", 42)
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("SetCancelOnDisconnect", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/set_cancel_on_disconnect", boom)
		_, err := api.SetCancelOnDisconnect(context.Background(), true)
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("ChangeSubaccountLabel", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/change_subaccount_label", boom)
		_, err := api.ChangeSubaccountLabel(context.Background(), "newlabel")
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetStatistics", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/statistics", boom)
		_, err := api.GetStatistics(context.Background(), "BTC-PERP")
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetTransaction", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/get_transaction", boom)
		_, err := api.GetTransaction(context.Background(), "TX1")
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetCurrencies", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/get_all_currencies", boom)
		_, err := api.GetCurrencies(context.Background())
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("CancelByLabel", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/cancel_by_label", boom)
		_, err := api.CancelByLabel(context.Background(), "L")
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("CancelByInstrument", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/cancel_by_instrument", boom)
		_, err := api.CancelByInstrument(context.Background(), "BTC-PERP")
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("CancelAll", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/cancel_all", boom)
		_, err := api.CancelAll(context.Background())
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetDepositHistory_ServerError", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_deposit_history", boom)
		_, _, err := api.GetDepositHistory(context.Background(), types.PageRequest{Page: 1, PageSize: 10})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetWithdrawalHistory_ServerError", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_withdrawal_history", boom)
		_, _, err := api.GetWithdrawalHistory(context.Background(), types.PageRequest{})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetTradeHistory_ServerError", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/get_trade_history", boom)
		_, _, err := api.GetTradeHistory(context.Background(), types.PageRequest{})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
}

// silence unused-import warning when none of the testutil symbols
// surface in this file (newAPI lives in the same package).
var _ = testutil.NewFakeTransport
