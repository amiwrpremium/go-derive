package methods_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

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

func TestGetOrderHistory(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_orders", map[string]any{
		"orders":     []any{},
		"pagination": map[string]any{"num_pages": 2, "count": 100},
	})
	_, page, err := api.GetOrderHistory(context.Background(), types.PageRequest{Page: 1, PageSize: 50})
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
			_, _, e := api.GetOrderHistory(context.Background(), types.PageRequest{})
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
