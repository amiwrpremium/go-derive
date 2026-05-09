package methods_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/methods"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func validPlaceOrderInput() methods.PlaceOrderInput {
	return methods.PlaceOrderInput{
		InstrumentName: "BTC-PERP",
		Asset:          common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		Amount:         types.MustDecimal("1"),
		LimitPrice:     types.MustDecimal("100"),
		MaxFee:         types.MustDecimal("1"),
	}
}

func TestPlaceOrderInput_Validate_Happy(t *testing.T) {
	require.NoError(t, validPlaceOrderInput().Validate())
}

func TestPlaceOrderInput_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*methods.PlaceOrderInput)
		want string
	}{
		{"empty instrument", func(in *methods.PlaceOrderInput) { in.InstrumentName = "" }, "instrument_name"},
		{"zero asset", func(in *methods.PlaceOrderInput) { in.Asset = common.Address{} }, "asset"},
		{"bad direction", func(in *methods.PlaceOrderInput) { in.Direction = enums.Direction("x") }, "direction"},
		{"bad order type", func(in *methods.PlaceOrderInput) { in.OrderType = enums.OrderType("x") }, "order_type"},
		{"bad time-in-force", func(in *methods.PlaceOrderInput) { in.TimeInForce = enums.TimeInForce("x") }, "time_in_force"},
		{"zero amount", func(in *methods.PlaceOrderInput) { in.Amount = types.MustDecimal("0") }, "amount"},
		{"zero price", func(in *methods.PlaceOrderInput) { in.LimitPrice = types.MustDecimal("0") }, "limit_price"},
		{"negative fee", func(in *methods.PlaceOrderInput) { in.MaxFee = types.MustDecimal("-1") }, "max_fee"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in := validPlaceOrderInput()
			c.mut(&in)
			err := in.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, types.ErrInvalidParams))
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
	_, err := api.PlaceOrder(context.Background(), methods.PlaceOrderInput{})
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
}

func TestPlaceOrder_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0) // signer set but subaccount=0
	_, err := api.PlaceOrder(context.Background(), methods.PlaceOrderInput{})
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestPlaceOrder_Success_PopulatesSignatureFields(t *testing.T) {
	api, ft := newAPI(t, true, 1)

	// Capture the params shipped on the wire.
	ft.Handle("private/order", func(_ json.RawMessage) (any, error) {
		// Echo the order id back in the canonical response shape.
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

	in := methods.PlaceOrderInput{
		InstrumentName: "BTC-PERP",
		Asset:          common.HexToAddress("0x1111111111111111111111111111111111111111"),
		SubID:          0,
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		TimeInForce:    enums.TimeInForceGTC,
		Amount:         types.MustDecimal("0.1"),
		LimitPrice:     types.MustDecimal("65000"),
		MaxFee:         types.MustDecimal("10"),
	}
	order, err := api.PlaceOrder(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, "O1", order.OrderID)

	// The captured params must include signature, signer, signature_expiry_sec, nonce.
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
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
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

func TestGetOrders_NoFilter(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_orders", map[string]any{
		"orders":     []any{},
		"pagination": map[string]any{"num_pages": 2, "count": 100},
	})
	_, page, err := api.GetOrders(context.Background(), types.PageRequest{Page: 1, PageSize: 50}, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, page.NumPages)
	assert.Equal(t, 100, page.Count)
	params := paramsAsMap(t, ft.LastCall().Params)
	_, hasInst := params["instrument_name"]
	_, hasLabel := params["label"]
	_, hasStatus := params["status"]
	assert.False(t, hasInst)
	assert.False(t, hasLabel)
	assert.False(t, hasStatus)
}

func TestGetOrders_WithFilter(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_orders", map[string]any{
		"orders":     []any{},
		"pagination": map[string]any{"num_pages": 1, "count": 0},
	})
	_, _, err := api.GetOrders(context.Background(), types.PageRequest{}, &methods.GetOrdersFilter{
		InstrumentName: "BTC-PERP",
		Label:          "alpha",
		Status:         enums.OrderStatusOpen,
	})
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "BTC-PERP", params["instrument_name"])
	assert.Equal(t, "alpha", params["label"])
	assert.Equal(t, "open", params["status"])
}

func TestGetOrders_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, _, err := api.GetOrders(context.Background(), types.PageRequest{}, nil)
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestGetOrderHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/get_order_history", map[string]any{
		"subaccount_id": int64(7),
		"orders": []any{
			map[string]any{
				"order_id": "O42", "subaccount_id": 7, "instrument_name": "BTC-PERP",
				"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
				"order_status": "filled", "amount": "0.1", "filled_amount": "0.1",
				"limit_price": "65000", "max_fee": "10", "nonce": 1,
				"signer":             "0x0000000000000000000000000000000000000000",
				"creation_timestamp": 1700000000000, "last_update_timestamp": 1700000060000,
			},
		},
		"pagination": map[string]any{"num_pages": 1, "count": 1},
	})
	orders, page, err := api.GetOrderHistory(context.Background(), types.PageRequest{}, methods.OrderHistoryQuery{
		FromTimestamp: types.MillisTime{T: time.UnixMilli(1700000000000)},
		ToTimestamp:   types.MillisTime{T: time.UnixMilli(1700000060000)},
	})
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, "O42", orders[0].OrderID)
	assert.Equal(t, 1, page.Count)
	// Subaccount must be threaded when Wallet is empty.
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, float64(7), params["subaccount_id"])
	assert.Equal(t, float64(1700000000000), params["from_timestamp"])
	assert.Equal(t, float64(1700000060000), params["to_timestamp"])
}

func TestGetOrderHistory_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, _, err := api.GetOrderHistory(context.Background(), types.PageRequest{}, methods.OrderHistoryQuery{})
	assert.True(t, errors.Is(err, derrors.ErrUnauthorized))
}

func TestGetOrderHistory_AcceptsWalletWithoutSubaccount(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/get_order_history", map[string]any{
		"subaccount_id": int64(0),
		"orders":        []any{},
		"pagination":    map[string]any{"num_pages": 0, "count": 0},
	})
	_, _, err := api.GetOrderHistory(context.Background(), types.PageRequest{}, methods.OrderHistoryQuery{Wallet: "0xabc"})
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	_, hasSub := params["subaccount_id"]
	assert.False(t, hasSub, "subaccount_id must not be set when wallet is provided")
	assert.Equal(t, "0xabc", params["wallet"])
}

func TestGetOrderHistory_RequiresSubaccountWhenWalletAbsent(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, _, err := api.GetOrderHistory(context.Background(), types.PageRequest{}, methods.OrderHistoryQuery{})
	assert.True(t, errors.Is(err, derrors.ErrSubaccountRequired))
}

func TestCancelTriggerOrder_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/cancel_trigger_order", map[string]any{
		"order_id": "T1", "subaccount_id": int64(7), "instrument_name": "BTC-PERP",
		"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
		"order_status": "cancelled", "amount": "0.1", "filled_amount": "0",
		"limit_price": "65000", "max_fee": "10", "nonce": int64(1),
		"signer":             "0x0000000000000000000000000000000000000000",
		"creation_timestamp": int64(1700000000000), "last_update_timestamp": int64(1700000060000),
	})
	got, err := api.CancelTriggerOrder(context.Background(), "T1")
	require.NoError(t, err)
	assert.Equal(t, "T1", got.OrderID)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "T1", params["order_id"])
	assert.Equal(t, float64(7), params["subaccount_id"])
}

func TestCancelTriggerOrder_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.CancelTriggerOrder(context.Background(), "T1")
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestCancelAllTriggerOrders(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/cancel_all_trigger_orders", "ok")
	require.NoError(t, api.CancelAllTriggerOrders(context.Background()))
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, float64(7), params["subaccount_id"])
}

func TestCancelAllTriggerOrders_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	err := api.CancelAllTriggerOrders(context.Background())
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestPrivateMethods_RequireSubaccount_Across(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	cases := map[string]func() error{
		"GetOrder":      func() error { _, e := api.GetOrder(context.Background(), "x"); return e },
		"GetOpenOrders": func() error { _, e := api.GetOpenOrders(context.Background()); return e },
		"GetOrders": func() error {
			_, _, e := api.GetOrders(context.Background(), types.PageRequest{}, nil)
			return e
		},
		"CancelByLabel":      func() error { _, e := api.CancelByLabel(context.Background(), "x"); return e },
		"CancelByInstrument": func() error { _, e := api.CancelByInstrument(context.Background(), "x"); return e },
		"CancelAll":          func() error { _, e := api.CancelAll(context.Background()); return e },
	}
	for name, fn := range cases {
		t.Run(name, func(t *testing.T) {
			assert.ErrorIs(t, fn(), derrors.ErrSubaccountRequired)
		})
	}
}

func TestReplace_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/replace", map[string]any{
		"cancelled_order": map[string]any{
			"order_id": "old", "subaccount_id": int64(1), "instrument_name": "BTC-PERP",
			"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
			"order_status": "cancelled", "amount": "1", "filled_amount": "0",
			"limit_price": "100", "max_fee": "5", "nonce": int64(1),
			"signer":                "0x0000000000000000000000000000000000000001",
			"creation_timestamp":    int64(1700000000000),
			"last_update_timestamp": int64(1700000001000),
		},
		"order": map[string]any{
			"order_id": "new", "subaccount_id": int64(1), "instrument_name": "BTC-PERP",
			"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
			"order_status": "open", "amount": "1", "filled_amount": "0",
			"limit_price": "101", "max_fee": "5", "nonce": int64(2),
			"signer":                "0x0000000000000000000000000000000000000001",
			"creation_timestamp":    int64(1700000002000),
			"last_update_timestamp": int64(1700000002000),
		},
		"trades": []any{},
	})
	res, err := api.Replace(context.Background(), map[string]any{"order_id_to_cancel": "old"})
	require.NoError(t, err)
	assert.Equal(t, "old", res.CancelledOrder.OrderID)
	require.NotNil(t, res.Order)
	assert.Equal(t, "new", res.Order.OrderID)
	assert.Empty(t, res.Trades)
}

func TestReplace_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleError("private/replace", boom)
	_, err := api.Replace(context.Background(), nil)
	assert.ErrorAs(t, err, new(*derrors.APIError))
}

func TestReplace_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.Replace(context.Background(), nil)
	assert.True(t, errors.Is(err, derrors.ErrUnauthorized))
}

func TestOrderDebug_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/order_debug", map[string]any{
		"action_hash":         "0xaa",
		"encoded_data":        "0xbb",
		"encoded_data_hashed": "0xcc",
		"typed_data_hash":     "0xdd",
		"raw_data": map[string]any{
			"data":              map[string]any{"asset": "0x1"},
			"expiry":            int64(1700000000),
			"is_atomic_signing": false,
			"module":            "0xmodule",
			"nonce":             int64(42),
			"owner":             "0xowner",
			"signature":         "0xsig",
			"signer":            "0xsigner",
			"subaccount_id":     int64(7),
		},
	})
	dbg, err := api.OrderDebug(context.Background(), map[string]any{"instrument_name": "BTC-PERP"})
	require.NoError(t, err)
	assert.Equal(t, "0xdd", dbg.TypedDataHash)
	assert.Equal(t, int64(42), dbg.RawData.Nonce)
	assert.Equal(t, int64(7), dbg.RawData.SubaccountID)
}

func TestOrderDebug_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleError("private/order_debug", boom)
	_, err := api.OrderDebug(context.Background(), nil)
	assert.ErrorAs(t, err, new(*derrors.APIError))
}

func TestCancelByNonce_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/cancel_by_nonce", map[string]any{"cancelled_orders": int64(2)})
	res, err := api.CancelByNonce(context.Background(), "BTC-PERP", 42)
	require.NoError(t, err)
	assert.Equal(t, int64(2), res.CancelledOrders)
}

func TestSetCancelOnDisconnect_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/set_cancel_on_disconnect", "ok")
	require.NoError(t, api.SetCancelOnDisconnect(context.Background(), true))
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, true, params["enabled"])
}

func TestSetCancelOnDisconnect_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	err := api.SetCancelOnDisconnect(context.Background(), true)
	assert.True(t, errors.Is(err, derrors.ErrUnauthorized))
}
