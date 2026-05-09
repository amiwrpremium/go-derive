package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestReplaceResult_HappyPath(t *testing.T) {
	raw := []byte(`{
		"cancelled_order":{"order_id":"old","subaccount_id":1,"instrument_name":"BTC-PERP","direction":"buy","order_type":"limit","time_in_force":"gtc","order_status":"cancelled","amount":"1","filled_amount":"0","limit_price":"100","max_fee":"5","nonce":1,"signer":"0x0000000000000000000000000000000000000001","creation_timestamp":1700000000000,"last_update_timestamp":1700000001000},
		"order":{"order_id":"new","subaccount_id":1,"instrument_name":"BTC-PERP","direction":"buy","order_type":"limit","time_in_force":"gtc","order_status":"open","amount":"1","filled_amount":"0","limit_price":"101","max_fee":"5","nonce":2,"signer":"0x0000000000000000000000000000000000000001","creation_timestamp":1700000002000,"last_update_timestamp":1700000002000},
		"trades":[]
	}`)
	var got types.ReplaceResult
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "old", got.CancelledOrder.OrderID)
	require.NotNil(t, got.Order)
	assert.Equal(t, "new", got.Order.OrderID)
	assert.Nil(t, got.CreateOrderError)
	assert.Empty(t, got.Trades)
}

func TestReplaceResult_CreateOrderError(t *testing.T) {
	raw := []byte(`{
		"cancelled_order":{"order_id":"old","subaccount_id":1,"instrument_name":"BTC-PERP","direction":"buy","order_type":"limit","time_in_force":"gtc","order_status":"cancelled","amount":"1","filled_amount":"0","limit_price":"100","max_fee":"5","nonce":1,"signer":"0x0000000000000000000000000000000000000001","creation_timestamp":1700000000000,"last_update_timestamp":1700000001000},
		"create_order_error":{"code":-32602,"message":"invalid","data":"limit_price"}
	}`)
	var got types.ReplaceResult
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Nil(t, got.Order)
	require.NotNil(t, got.CreateOrderError)
	assert.Equal(t, -32602, got.CreateOrderError.Code)
	assert.Equal(t, "invalid", got.CreateOrderError.Message)
	assert.Equal(t, "limit_price", got.CreateOrderError.Data)
}

func TestOrderDebugResult_Decode(t *testing.T) {
	raw := []byte(`{
		"action_hash":"0xaaa",
		"encoded_data":"0xbbb",
		"encoded_data_hashed":"0xccc",
		"typed_data_hash":"0xddd",
		"raw_data":{
			"data":{"asset":"0x123"},
			"expiry":1700000000,
			"is_atomic_signing":false,
			"module":"0xmodule",
			"nonce":42,
			"owner":"0xowner",
			"signature":"0xsig",
			"signer":"0xsigner",
			"subaccount_id":7
		}
	}`)
	var got types.OrderDebugResult
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "0xaaa", got.ActionHash)
	assert.Equal(t, "0xddd", got.TypedDataHash)
	assert.Equal(t, int64(1700000000), got.RawData.Expiry)
	assert.Equal(t, "0xmodule", got.RawData.Module)
	assert.Equal(t, int64(42), got.RawData.Nonce)
	assert.Equal(t, int64(7), got.RawData.SubaccountID)
	assert.JSONEq(t, `{"asset":"0x123"}`, string(got.RawData.Data))
}

func TestCancelByNonceResult_Decode(t *testing.T) {
	raw := []byte(`{"cancelled_orders":3}`)
	var got types.CancelByNonceResult
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, int64(3), got.CancelledOrders)
}
