package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestGetTradeHistory_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_trade_history", map[string]any{
		"trades":     []any{},
		"pagination": map[string]any{"num_pages": 4, "count": 100},
	})
	_, page, err := api.GetTradeHistory(context.Background(), types.PageRequest{Page: 1, PageSize: 25})
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
	_, _, err := api.GetTradeHistory(context.Background(), types.PageRequest{})
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	_, hasPage := params["page"]
	_, hasSize := params["page_size"]
	assert.False(t, hasPage)
	assert.False(t, hasSize)
}

func TestGetTradeHistory_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, _, err := api.GetTradeHistory(context.Background(), types.PageRequest{})
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestGetDepositHistory_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_deposit_history", map[string]any{
		"events":     []any{},
		"pagination": map[string]any{"num_pages": 0, "count": 0, "current_page": 1, "page_size": 10},
	})
	_, _, err := api.GetDepositHistory(context.Background(), types.PageRequest{Page: 0, PageSize: 0})
	require.NoError(t, err)
}

func TestGetDepositHistory_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, _, err := api.GetDepositHistory(context.Background(), types.PageRequest{})
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestGetWithdrawalHistory_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_withdrawal_history", map[string]any{
		"events":     []any{},
		"pagination": map[string]any{"num_pages": 0, "count": 0, "current_page": 1, "page_size": 10},
	})
	_, _, err := api.GetWithdrawalHistory(context.Background(), types.PageRequest{Page: 2, PageSize: 50})
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, float64(2), params["page"])
	assert.Equal(t, float64(50), params["page_size"])
}

func TestGetWithdrawalHistory_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, _, err := api.GetWithdrawalHistory(context.Background(), types.PageRequest{})
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestGetTransaction_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_transaction", map[string]any{
		"data":             "{\"foo\":\"bar\"}",
		"error_log":        nil,
		"status":           "settled",
		"transaction_hash": "0xabc",
	})
	got, err := api.GetTransaction(context.Background(), "tx-1")
	require.NoError(t, err)
	assert.Equal(t, "settled", got.Status)
	assert.Equal(t, "0xabc", got.TransactionHash)
	assert.Equal(t, "", got.ErrorLog)
}

func TestGetTransaction_FailedTx(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_transaction", map[string]any{
		"data":             "...",
		"error_log":        "reverted: insufficient",
		"status":           "reverted",
		"transaction_hash": nil,
	})
	got, err := api.GetTransaction(context.Background(), "tx-1")
	require.NoError(t, err)
	assert.Equal(t, "reverted", got.Status)
	assert.Equal(t, "reverted: insufficient", got.ErrorLog)
}
