package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
)

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
