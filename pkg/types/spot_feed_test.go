package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestSpotFeed_Decode(t *testing.T) {
	// Mirrors a real `spot_feed.BTC` notification payload.
	raw := `{
		"timestamp": 1777842232556,
		"feeds": {
			"BTC": {
				"price": "78908.29",
				"confidence": "1",
				"price_prev_daily": "78689.04",
				"confidence_prev_daily": "1",
				"timestamp_prev_daily": 1777755832556
			}
		}
	}`
	var sf types.SpotFeed
	require.NoError(t, json.Unmarshal([]byte(raw), &sf))
	assert.Equal(t, int64(1777842232556), sf.Timestamp.Millis())
	require.Contains(t, sf.Feeds, "BTC")
	btc := sf.Feeds["BTC"]
	assert.Equal(t, "78908.29", btc.Price.String())
	assert.Equal(t, "78689.04", btc.PricePrevDaily.String())
	assert.Equal(t, int64(1777755832556), btc.TimestampPrevDaily.Millis())
}
