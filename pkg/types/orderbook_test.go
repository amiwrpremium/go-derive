package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestOrderBookLevel_JSONRoundTrip(t *testing.T) {
	in := types.OrderBookLevel{
		Price:  types.MustDecimal("100.5"),
		Amount: types.MustDecimal("0.25"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `["100.5","0.25"]`, string(b))

	var out types.OrderBookLevel
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.Price.String(), out.Price.String())
	assert.Equal(t, in.Amount.String(), out.Amount.String())
}

func TestOrderBookLevel_UnmarshalRejectsObject(t *testing.T) {
	var l types.OrderBookLevel
	err := json.Unmarshal([]byte(`{"price":"1","amount":"2"}`), &l)
	assert.Error(t, err)
}

func TestOrderBook_DecodeFullPayload(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"bids": [["65000","0.1"],["64999","0.2"]],
		"asks": [["65001","0.3"]],
		"timestamp": 1700000000000,
		"publish_time": 1700000000005
	}`
	var ob types.OrderBook
	require.NoError(t, json.Unmarshal([]byte(payload), &ob))
	assert.Equal(t, "BTC-PERP", ob.InstrumentName)
	require.Len(t, ob.Bids, 2)
	require.Len(t, ob.Asks, 1)
	assert.Equal(t, "65000", ob.Bids[0].Price.String())
	assert.Equal(t, "0.1", ob.Bids[0].Amount.String())
	assert.Equal(t, "65001", ob.Asks[0].Price.String())
	assert.Equal(t, int64(1700000000000), ob.Timestamp.Millis())
}
