package methods_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

const instAddrHex = "0x1111111111111111111111111111111111111111"

// methodsOf extracts just the method names from a transport's call log.
// Tests use it to assert the SDK fanned a call out to (or away from)
// the expected upstream methods.
func methodsOf(calls []testutil.FakeCall) []string {
	out := make([]string, len(calls))
	for i, c := range calls {
		out[i] = c.Method
	}
	return out
}

func instrumentResponse(name, subID string) map[string]any {
	return map[string]any{
		"instrument_name":   name,
		"base_currency":     "BTC",
		"quote_currency":    "USDC",
		"instrument_type":   "perp",
		"is_active":         true,
		"tick_size":         "0.5",
		"minimum_amount":    "0.001",
		"maximum_amount":    "1000",
		"amount_step":       "0.001",
		"base_asset_address":  instAddrHex,
		"base_asset_sub_id":   subID,
	}
}

func TestPlaceOrder_AutoResolvesAssetFromInstrumentName(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("public/get_instrument", instrumentResponse("BTC-PERP", "42"))
	ft.HandleResult("private/order", map[string]any{
		"order": map[string]any{
			"order_id": "O1", "subaccount_id": 1, "instrument_name": "BTC-PERP",
			"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
			"order_status": "open", "amount": "0.1", "filled_amount": "0",
			"limit_price": "65000", "max_fee": "10", "nonce": 1,
			"signer":             "0x0000000000000000000000000000000000000000",
			"creation_timestamp": 1, "last_update_timestamp": 1,
		},
	})

	// Caller leaves Asset/SubID zero; the SDK should fetch them via
	// get_instrument before signing.
	_, _, err := api.PlaceOrder(context.Background(), types.PlaceOrderInput{
		InstrumentName: "BTC-PERP",
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		Amount:         types.MustDecimal("0.1"),
		LimitPrice:     types.MustDecimal("65000"),
		MaxFee:         types.MustDecimal("10"),
	})
	require.NoError(t, err)
	methods := methodsOf(ft.Calls())
	assert.Equal(t, []string{"public/get_instrument", "private/order"}, methods,
		"expected one get_instrument lookup before private/order")
}

func TestPlaceOrder_CachesInstrumentAcrossCalls(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("public/get_instrument", instrumentResponse("BTC-PERP", "0"))
	ft.HandleResult("private/order", map[string]any{
		"order": map[string]any{
			"order_id": "O1", "subaccount_id": 1, "instrument_name": "BTC-PERP",
			"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
			"order_status": "open", "amount": "0.1", "filled_amount": "0",
			"limit_price": "65000", "max_fee": "10", "nonce": 1,
			"signer":             "0x0000000000000000000000000000000000000000",
			"creation_timestamp": 1, "last_update_timestamp": 1,
		},
	})

	in := types.PlaceOrderInput{
		InstrumentName: "BTC-PERP",
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		Amount:         types.MustDecimal("0.1"),
		LimitPrice:     types.MustDecimal("65000"),
		MaxFee:         types.MustDecimal("10"),
	}
	_, _, err := api.PlaceOrder(context.Background(), in)
	require.NoError(t, err)
	_, _, err = api.PlaceOrder(context.Background(), in)
	require.NoError(t, err)

	var instrumentCalls, orderCalls int
	for _, m := range methodsOf(ft.Calls()) {
		switch m {
		case "public/get_instrument":
			instrumentCalls++
		case "private/order":
			orderCalls++
		}
	}
	assert.Equal(t, 1, instrumentCalls, "second PlaceOrder should hit the cache")
	assert.Equal(t, 2, orderCalls)
}

func TestPlaceOrder_PreservesExplicitAssetAndSkipsLookup(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/order", map[string]any{
		"order": map[string]any{
			"order_id": "O1", "subaccount_id": 1, "instrument_name": "BTC-PERP",
			"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
			"order_status": "open", "amount": "0.1", "filled_amount": "0",
			"limit_price": "65000", "max_fee": "10", "nonce": 1,
			"signer":             "0x0000000000000000000000000000000000000000",
			"creation_timestamp": 1, "last_update_timestamp": 1,
		},
	})

	_, _, err := api.PlaceOrder(context.Background(), types.PlaceOrderInput{
		InstrumentName: "BTC-PERP",
		Asset:          types.Address(common.HexToAddress(instAddrHex)),
		SubID:          99,
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		Amount:         types.MustDecimal("0.1"),
		LimitPrice:     types.MustDecimal("65000"),
		MaxFee:         types.MustDecimal("10"),
	})
	require.NoError(t, err)
	for _, m := range methodsOf(ft.Calls()) {
		assert.NotEqual(t, "public/get_instrument", m,
			"explicit Asset must bypass the cache lookup entirely")
	}
}

func TestGetInstruments_PopulatesCacheForSubsequentPlaceOrder(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("public/get_instruments", []map[string]any{
		instrumentResponse("BTC-PERP", "0"),
		instrumentResponse("ETH-PERP", "0"),
	})
	ft.HandleResult("private/order", map[string]any{
		"order": map[string]any{
			"order_id": "O1", "subaccount_id": 1, "instrument_name": "BTC-PERP",
			"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
			"order_status": "open", "amount": "0.1", "filled_amount": "0",
			"limit_price": "65000", "max_fee": "10", "nonce": 1,
			"signer":             "0x0000000000000000000000000000000000000000",
			"creation_timestamp": 1, "last_update_timestamp": 1,
		},
	})

	insts, err := api.GetInstruments(context.Background(), "BTC", enums.InstrumentTypePerp)
	require.NoError(t, err)
	require.Len(t, insts, 2)

	// Now PlaceOrder should NOT trigger a get_instrument call.
	_, _, err = api.PlaceOrder(context.Background(), types.PlaceOrderInput{
		InstrumentName: "BTC-PERP",
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		Amount:         types.MustDecimal("0.1"),
		LimitPrice:     types.MustDecimal("65000"),
		MaxFee:         types.MustDecimal("10"),
	})
	require.NoError(t, err)
	for _, m := range methodsOf(ft.Calls()) {
		assert.NotEqual(t, "public/get_instrument", m,
			"GetInstruments must populate the cache as a side effect")
	}
}

func TestResolveInstrument_FailsWhenInstrumentLacksOnChainMetadata(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	// Instrument returned with empty base_asset_address — can't sign.
	ft.HandleResult("public/get_instrument", map[string]any{
		"instrument_name":    "GHOST",
		"base_currency":      "BTC",
		"quote_currency":     "USDC",
		"instrument_type":    "perp",
		"is_active":          true,
		"tick_size":          "0.5",
		"minimum_amount":     "0.001",
		"maximum_amount":     "1000",
		"amount_step":        "0.001",
		"base_asset_address": "0x0000000000000000000000000000000000000000",
		"base_asset_sub_id":  "0",
	})

	_, _, err := api.PlaceOrder(context.Background(), types.PlaceOrderInput{
		InstrumentName: "GHOST",
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		Amount:         types.MustDecimal("0.1"),
		LimitPrice:     types.MustDecimal("65000"),
		MaxFee:         types.MustDecimal("10"),
	})
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "GHOST"))
}

func TestInvalidateInstrumentCache_ForcesRefetch(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("public/get_instrument", instrumentResponse("BTC-PERP", "0"))
	ft.HandleResult("private/order", map[string]any{
		"order": map[string]any{
			"order_id": "O1", "subaccount_id": 1, "instrument_name": "BTC-PERP",
			"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
			"order_status": "open", "amount": "0.1", "filled_amount": "0",
			"limit_price": "65000", "max_fee": "10", "nonce": 1,
			"signer":             "0x0000000000000000000000000000000000000000",
			"creation_timestamp": 1, "last_update_timestamp": 1,
		},
	})

	in := types.PlaceOrderInput{
		InstrumentName: "BTC-PERP",
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		Amount:         types.MustDecimal("0.1"),
		LimitPrice:     types.MustDecimal("65000"),
		MaxFee:         types.MustDecimal("10"),
	}
	_, _, err := api.PlaceOrder(context.Background(), in)
	require.NoError(t, err)
	api.InvalidateInstrumentCache("BTC-PERP")
	_, _, err = api.PlaceOrder(context.Background(), in)
	require.NoError(t, err)

	var instrumentCalls int
	for _, m := range methodsOf(ft.Calls()) {
		if m == "public/get_instrument" {
			instrumentCalls++
		}
	}
	assert.Equal(t, 2, instrumentCalls,
		"InvalidateInstrumentCache must drop the entry so the next call refetches")
}

func TestPreloadInstruments_FetchesPerCurrencyAndKind(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("public/get_instruments", []map[string]any{
		instrumentResponse("BTC-PERP", "0"),
	})
	err := api.PreloadInstruments(context.Background(), "BTC", "ETH")
	require.NoError(t, err)
	// 2 currencies × 3 kinds = 6 calls.
	calls := 0
	for _, m := range methodsOf(ft.Calls()) {
		if m == "public/get_instruments" {
			calls++
		}
	}
	assert.Equal(t, 6, calls)
}

func TestPreloadInstruments_PropagatesError(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleError("public/get_instruments", errors.New("upstream"))
	err := api.PreloadInstruments(context.Background(), "BTC")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "BTC")
}
