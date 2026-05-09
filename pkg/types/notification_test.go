package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestNotification_Decode(t *testing.T) {
	raw := []byte(`{
		"id":42,
		"subaccount_id":7,
		"event":"deposit_completed",
		"event_details":{"asset":"USDC","amount":"100"},
		"status":"unseen",
		"timestamp":1700000000000,
		"transaction_id":1234,
		"tx_hash":"0xabc"
	}`)
	var got types.Notification
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, int64(42), got.ID)
	assert.Equal(t, int64(7), got.SubaccountID)
	assert.Equal(t, "deposit_completed", got.Event)
	assert.JSONEq(t, `{"asset":"USDC","amount":"100"}`, string(got.EventDetails))
	assert.Equal(t, "unseen", got.Status)
	assert.Equal(t, int64(1700000000000), got.Timestamp.Millis())
	require.NotNil(t, got.TransactionID)
	assert.Equal(t, int64(1234), *got.TransactionID)
	assert.Equal(t, "0xabc", got.TxHash)
}

func TestNotification_NullableFields(t *testing.T) {
	// transaction_id and tx_hash are nullable on the wire.
	raw := []byte(`{
		"id":1,
		"subaccount_id":1,
		"event":"info",
		"event_details":{},
		"status":"seen",
		"timestamp":1700000000000,
		"transaction_id":null,
		"tx_hash":null
	}`)
	var got types.Notification
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Nil(t, got.TransactionID)
	assert.Equal(t, "", got.TxHash)
}

func TestUpdateNotificationsResult_Decode(t *testing.T) {
	raw := []byte(`{"updated_count":5}`)
	var got types.UpdateNotificationsResult
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, int64(5), got.UpdatedCount)
}
