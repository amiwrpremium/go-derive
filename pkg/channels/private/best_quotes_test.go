package private_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/channels/private"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestBestQuotes_Name(t *testing.T) {
	got := private.BestQuotes{SubaccountID: 42}.Name()
	assert.Equal(t, "42.best.quotes", got)
}

func TestBestQuotes_Decode_Result(t *testing.T) {
	raw := json.RawMessage(`[{
		"rfq_id":"R1",
		"error":null,
		"result":{
			"direction":"buy",
			"is_valid":true,
			"invalid_reason":null,
			"estimated_fee":"5",
			"estimated_realized_pnl":"0",
			"estimated_realized_pnl_excl_fees":"0",
			"estimated_total_cost":"100",
			"filled_pct":"0",
			"orderbook_total_cost":"100",
			"suggested_max_fee":"10",
			"pre_initial_margin":"100",
			"post_initial_margin":"110",
			"post_liquidation_price":"45000",
			"down_liquidation_price":"40000",
			"up_liquidation_price":"60000"
		}
	}]`)
	v, err := private.BestQuotes{}.Decode(raw)
	require.NoError(t, err)
	events, ok := v.([]types.BestQuoteFeedEvent)
	require.True(t, ok, "Decode must return []types.BestQuoteFeedEvent")
	require.Len(t, events, 1)
	assert.Equal(t, "R1", events[0].RFQID)
	assert.Nil(t, events[0].Error)
	require.NotNil(t, events[0].Result)
	assert.True(t, events[0].Result.IsValid)
}

func TestBestQuotes_Decode_Error(t *testing.T) {
	raw := json.RawMessage(`[{
		"rfq_id":"R1",
		"error":{"code":-32000,"message":"insufficient_margin"},
		"result":null
	}]`)
	v, err := private.BestQuotes{}.Decode(raw)
	require.NoError(t, err)
	events, ok := v.([]types.BestQuoteFeedEvent)
	require.True(t, ok)
	require.Len(t, events, 1)
	require.NotNil(t, events[0].Error)
	assert.Equal(t, -32000, events[0].Error.Code)
	assert.Nil(t, events[0].Result)
}

func TestBestQuotes_Decode_Malformed(t *testing.T) {
	_, err := private.BestQuotes{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
