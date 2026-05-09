package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestSpotFeedCandle_Decode(t *testing.T) {
	raw := []byte(`{
		"timestamp": 1700000000000,
		"timestamp_bucket": 1700000000000,
		"price": "2500",
		"open_price": "2495",
		"high_price": "2510",
		"low_price": "2490",
		"close_price": "2500"
	}`)
	var c types.SpotFeedCandle
	require.NoError(t, json.Unmarshal(raw, &c))
	assert.Equal(t, "2500", c.Price.String())
	assert.Equal(t, "2495", c.OpenPrice.String())
	assert.Equal(t, "2510", c.HighPrice.String())
	assert.Equal(t, int64(1700000000000), c.Timestamp.Time().UnixMilli())
}

func TestTradingViewChart_Decode(t *testing.T) {
	raw := []byte(`{
		"timestamp": 1700000000000,
		"timestamp_bucket": 1700000000000,
		"open_price": "65000",
		"high_price": "65100",
		"low_price": "64900",
		"close_price": "65050",
		"volume_contracts": "10",
		"volume_usd": "650500"
	}`)
	var c types.TradingViewChart
	require.NoError(t, json.Unmarshal(raw, &c))
	assert.Equal(t, "65000", c.OpenPrice.String())
	assert.Equal(t, "10", c.VolumeContracts.String())
	assert.Equal(t, "650500", c.VolumeUSD.String())
}
