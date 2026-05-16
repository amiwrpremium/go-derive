package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
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
	assert.Equal(t, enums.DirectionBuy, o.Direction)
	assert.Equal(t, enums.OrderStatusOpen, o.OrderStatus)
	assert.True(t, o.ReduceOnly)
	assert.Equal(t, "alpha", o.Label)
}
