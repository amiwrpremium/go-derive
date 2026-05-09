package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestOrder_DecodeFull(t *testing.T) {
	payload := `{
		"order_id": "O1",
		"subaccount_id": 1,
		"instrument_name": "BTC-PERP",
		"direction": "buy",
		"order_type": "limit",
		"time_in_force": "gtc",
		"order_status": "open",
		"amount": "0.1",
		"filled_amount": "0",
		"limit_price": "65000",
		"average_price": "0",
		"max_fee": "10",
		"nonce": 1,
		"signer": "0x1111111111111111111111111111111111111111",
		"label": "alpha",
		"cancel_reason": "",
		"mmp": false,
		"reduce_only": true,
		"creation_timestamp": 1700000000000,
		"last_update_timestamp": 1700000000005
	}`
	var o types.Order
	require.NoError(t, json.Unmarshal([]byte(payload), &o))
	assert.Equal(t, "O1", o.OrderID)
	assert.Equal(t, derive.DirectionBuy, o.Direction)
	assert.Equal(t, derive.OrderStatusOpen, o.OrderStatus)
	assert.True(t, o.ReduceOnly)
	assert.Equal(t, "alpha", o.Label)
}

func TestOrderParams_OmitsEmptyOptionalFields(t *testing.T) {
	in := types.OrderParams{
		InstrumentName:  "BTC-PERP",
		Direction:       derive.DirectionBuy,
		OrderType:       derive.OrderTypeLimit,
		Amount:          types.MustDecimal("1"),
		LimitPrice:      types.MustDecimal("100"),
		MaxFee:          types.MustDecimal("10"),
		SubaccountID:    1,
		Nonce:           1,
		Signer:          types.MustAddress("0x1111111111111111111111111111111111111111"),
		Signature:       "0xdeadbeef",
		SignatureExpiry: 1700000000,
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	s := string(b)
	assert.NotContains(t, s, `"label"`)
	assert.NotContains(t, s, `"trigger_type"`)
	assert.NotContains(t, s, `"reduce_only":true`)
}

func TestOrderParams_IncludesPopulatedOptionals(t *testing.T) {
	in := types.OrderParams{
		InstrumentName: "BTC-PERP",
		Direction:      derive.DirectionSell,
		OrderType:      derive.OrderTypeLimit,
		Label:          "lbl",
		ReduceOnly:     true,
		MMP:            true,
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	s := string(b)
	assert.Contains(t, s, `"label":"lbl"`)
	assert.Contains(t, s, `"reduce_only":true`)
	assert.Contains(t, s, `"mmp":true`)
}

func TestCancelOrderParams_RoundTrip(t *testing.T) {
	in := types.CancelOrderParams{
		SubaccountID:   1,
		InstrumentName: "BTC-PERP",
		OrderID:        "O1",
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.CancelOrderParams
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}

func TestReplaceOrderParams_Embedded(t *testing.T) {
	in := types.ReplaceOrderParams{
		OrderIDToCancel: "O1",
		NewOrder: types.OrderParams{
			InstrumentName: "BTC-PERP",
			Direction:      derive.DirectionBuy,
		},
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"order_id_to_cancel":"O1"`)
	assert.Contains(t, string(b), `"new_order"`)
}
