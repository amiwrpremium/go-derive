package private_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/channels/private"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestOrders_Name(t *testing.T) {
	assert.Equal(t, "subaccount.123.orders", private.Orders{SubaccountID: 123}.Name())
}

func TestOrders_Name_Zero(t *testing.T) {
	assert.Equal(t, "subaccount.0.orders", private.Orders{}.Name())
}

func TestOrders_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"order_id":"O1","subaccount_id":1,"instrument_name":"BTC-PERP","direction":"buy","order_type":"limit","time_in_force":"gtc","order_status":"open","amount":"0.1","filled_amount":"0","limit_price":"65000","max_fee":"10","nonce":1,"signer":"0x0000000000000000000000000000000000000000","creation_timestamp":1700000000000,"last_update_timestamp":1700000000000}]`)
	v, err := private.Orders{}.Decode(raw)
	require.NoError(t, err)
	orders, ok := v.([]types.Order)
	require.True(t, ok)
	require.Len(t, orders, 1)
	assert.Equal(t, "O1", orders[0].OrderID)
}

func TestOrders_Decode_Malformed(t *testing.T) {
	_, err := private.Orders{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
