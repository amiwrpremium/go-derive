package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestLiquidationAuction_Decode(t *testing.T) {
	payload := `{
		"auction_id": "auc-9",
		"auction_type": "solvent",
		"bids": [],
		"end_timestamp": 1700000060000,
		"fee": "1.5",
		"start_timestamp": 1700000000000,
		"subaccount_id": 9,
		"tx_hash": "0x1111111111111111111111111111111111111111111111111111111111111111"
	}`
	var l types.LiquidationAuction
	require.NoError(t, json.Unmarshal([]byte(payload), &l))
	assert.Equal(t, "auc-9", l.AuctionID)
	assert.Equal(t, enums.AuctionTypeSolvent, l.AuctionType)
	assert.Equal(t, int64(9), l.SubaccountID)
	assert.Equal(t, "1.5", l.Fee.String())
	assert.False(t, l.TxHash.IsZero())
	assert.Equal(t, int64(1700000000000), l.StartTimestamp.Millis())
	assert.Equal(t, int64(1700000060000), l.EndTimestamp.Millis())
}

func TestLiquidationAuction_NullEndTimestamp(t *testing.T) {
	// The `end_timestamp` field is nullable on the wire — an open
	// auction reports `null`. Decode must leave the field zero.
	payload := `{
		"auction_id":"auc-1",
		"auction_type":"insolvent",
		"bids":[],
		"end_timestamp":null,
		"fee":"0",
		"start_timestamp":1700000000000,
		"subaccount_id":42,
		"tx_hash":"0x0000000000000000000000000000000000000000000000000000000000000000"
	}`
	var l types.LiquidationAuction
	require.NoError(t, json.Unmarshal([]byte(payload), &l))
	assert.Equal(t, enums.AuctionTypeInsolvent, l.AuctionType)
	assert.True(t, l.EndTimestamp.Time().IsZero(), "null end_timestamp should leave the field zero")
}

// Custom MarshalJSON on TxHash defeats omitempty. The wire format
// always includes tx_hash even when zero — document and pin that.
func TestLiquidationAuction_ZeroTxHashSerializedAsZeroString(t *testing.T) {
	in := types.LiquidationAuction{SubaccountID: 1}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"tx_hash":"0x0000000000000000000000000000000000000000000000000000000000000000"`)
}
