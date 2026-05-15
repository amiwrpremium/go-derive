package methods_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// pagedOrderHandler returns a fake "private/get_order_history" handler
// that emits a fresh page with one order each time it's called, until
// `pages` calls have been served. Each page reports NumPages=pages.
func pagedOrderHandler(t *testing.T, pages int) func(json.RawMessage) (any, error) {
	t.Helper()
	return func(raw json.RawMessage) (any, error) {
		var p struct {
			Page int `json:"page"`
		}
		require.NoError(t, json.Unmarshal(raw, &p))
		if p.Page > pages {
			return map[string]any{"orders": []any{}, "pagination": map[string]any{"num_pages": pages, "count": pages}}, nil
		}
		return map[string]any{
			"orders": []any{
				map[string]any{
					"order_id": "O", "subaccount_id": 1, "instrument_name": "BTC-PERP",
					"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
					"order_status": "open", "amount": "0.1", "filled_amount": "0",
					"limit_price": "65000", "max_fee": "10", "nonce": 1,
					"signer":             "0x0000000000000000000000000000000000000000",
					"creation_timestamp": 1, "last_update_timestamp": 1,
				},
			},
			"pagination": map[string]any{"num_pages": pages, "count": pages},
		}, nil
	}
}

func TestGetOrderHistoryAll_AccumulatesPagesAndStops(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.Handle("private/get_order_history", pagedOrderHandler(t, 3))

	all, err := api.GetOrderHistoryAll(context.Background(), types.OrderHistoryQuery{}, types.PaginateOptions{})
	require.NoError(t, err)
	assert.Len(t, all, 3, "should accumulate one order per page across all 3 pages")
}

func TestGetOrderHistoryAll_HonoursMaxItems(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.Handle("private/get_order_history", pagedOrderHandler(t, 5))

	all, err := api.GetOrderHistoryAll(context.Background(), types.OrderHistoryQuery{},
		types.PaginateOptions{MaxItems: 2})
	require.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestGetOrdersAll_ThreadsFilter(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	calls := 0
	ft.Handle("private/get_orders", func(raw json.RawMessage) (any, error) {
		calls++
		var p map[string]any
		require.NoError(t, json.Unmarshal(raw, &p))
		assert.Equal(t, "BTC-PERP", p["instrument_name"], "filter must reach the wire on every page")
		if calls == 1 {
			return map[string]any{
				"orders":     []any{map[string]any{"order_id": "O1", "subaccount_id": 1, "instrument_name": "BTC-PERP", "direction": "buy", "order_type": "limit", "time_in_force": "gtc", "order_status": "open", "amount": "0.1", "filled_amount": "0", "limit_price": "65000", "max_fee": "10", "nonce": 1, "signer": "0x0000000000000000000000000000000000000000", "creation_timestamp": 1, "last_update_timestamp": 1}},
				"pagination": map[string]any{"num_pages": 2, "count": 2},
			}, nil
		}
		return map[string]any{
			"orders":     []any{map[string]any{"order_id": "O2", "subaccount_id": 1, "instrument_name": "BTC-PERP", "direction": "buy", "order_type": "limit", "time_in_force": "gtc", "order_status": "open", "amount": "0.1", "filled_amount": "0", "limit_price": "65000", "max_fee": "10", "nonce": 1, "signer": "0x0000000000000000000000000000000000000000", "creation_timestamp": 1, "last_update_timestamp": 1}},
			"pagination": map[string]any{"num_pages": 2, "count": 2},
		}, nil
	})

	all, err := api.GetOrdersAll(context.Background(),
		&types.GetOrdersFilter{InstrumentName: "BTC-PERP"},
		types.PaginateOptions{})
	require.NoError(t, err)
	assert.Len(t, all, 2)
	assert.Equal(t, 2, calls)
}

func TestGetDepositHistoryAll_NoQueryArg(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_deposit_history", map[string]any{
		"events":     []any{},
		"pagination": map[string]any{"num_pages": 0, "count": 0},
	})
	all, err := api.GetDepositHistoryAll(context.Background(), types.PaginateOptions{})
	require.NoError(t, err)
	assert.Empty(t, all)
}

func TestGetAllInstrumentsAll_ThreadsKindAndExpired(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	calls := 0
	ft.Handle("public/get_all_instruments", func(raw json.RawMessage) (any, error) {
		calls++
		var p map[string]any
		require.NoError(t, json.Unmarshal(raw, &p))
		assert.Equal(t, "option", p["instrument_type"])
		assert.Equal(t, true, p["expired"], "includeExpired must reach the wire on every page")
		return map[string]any{
			"instruments": []any{},
			"pagination":  map[string]any{"num_pages": 0, "count": 0},
		}, nil
	})

	all, err := api.GetAllInstrumentsAll(context.Background(),
		types.AllInstrumentsQuery{Kind: enums.InstrumentTypeOption, IncludeExpired: true},
		types.PaginateOptions{})
	require.NoError(t, err)
	assert.Empty(t, all)
	assert.Equal(t, 1, calls)
}
