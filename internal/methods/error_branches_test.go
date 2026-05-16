package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/internal/testutil"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// boom is the error injected by the fake transport on every wrapper
// under test below. Using a single sentinel keeps the table compact.
var boom = &derrors.APIError{Code: 9999, Message: "boom"}

// TestNoArgWrappers_PropagateError exercises the APIError-pass-
// through path on every method, including the formerly-raw
// `GetOrderQuote` and `GetPerpImpactTWAP` (now both fully typed
// against `derivexyz/cockpit/orderbook-types`).
func TestNoArgWrappers_PropagateError(t *testing.T) {
	t.Run("CancelByNonce", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/cancel_by_nonce", boom)
		_, err := api.CancelByNonce(context.Background(), types.CancelByNonceInput{InstrumentName: "BTC-PERP", Nonce: 42})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("SetCancelOnDisconnect", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/set_cancel_on_disconnect", boom)
		err := api.SetCancelOnDisconnect(context.Background(), types.SetCancelOnDisconnectInput{Enabled: true})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("ChangeSubaccountLabel", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/change_subaccount_label", boom)
		err := api.ChangeSubaccountLabel(context.Background(), types.ChangeSubaccountLabelInput{Label: "newlabel"})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetStatistics", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/statistics", boom)
		_, err := api.GetStatistics(context.Background(), types.StatisticsQuery{InstrumentName: "BTC-PERP"})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetTransaction", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/get_transaction", boom)
		_, err := api.GetTransaction(context.Background(), types.TransactionQuery{TransactionID: "TX1"})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetFundingRateHistory", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/get_funding_rate_history", boom)
		_, err := api.GetFundingRateHistory(context.Background(), types.FundingRateHistoryQuery{InstrumentName: "BTC-PERP"})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetSpotFeedHistory", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/get_spot_feed_history", boom)
		_, _, err := api.GetSpotFeedHistory(context.Background(), types.SpotFeedHistoryQuery{Currency: "BTC", PeriodSec: 60})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetLatestSignedFeeds", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/get_latest_signed_feeds", boom)
		_, err := api.GetLatestSignedFeeds(context.Background(), types.LatestSignedFeedsQuery{})
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
		_, err := api.CancelByLabel(context.Background(), types.CancelByLabelInput{Label: "L"})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("CancelByInstrument", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/cancel_by_instrument", boom)
		_, err := api.CancelByInstrument(context.Background(), types.CancelByInstrumentInput{InstrumentName: "BTC-PERP"})
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
		_, _, err := api.GetTradeHistory(context.Background(), types.TradeHistoryQuery{}, types.PageRequest{})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetOrderQuote", func(t *testing.T) {
		api, ft := newAPI(t, true, 1)
		ft.HandleError("private/order_quote", boom)
		_, err := api.GetOrderQuote(context.Background(), validPlaceOrderInput())
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
	t.Run("GetPerpImpactTWAP", func(t *testing.T) {
		api, ft := newAPI(t, true, 0)
		ft.HandleError("public/get_perp_impact_twap", boom)
		_, err := api.GetPerpImpactTWAP(context.Background(), types.PerpImpactTWAPQuery{Currency: "BTC", EndTime: 1})
		assert.ErrorAs(t, err, new(*derrors.APIError))
	})
}

// silence unused-import warning when none of the testutil symbols
// surface in this file (newAPI lives in the same package).
var _ = testutil.NewFakeTransport
