package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestStatistics_Decode(t *testing.T) {
	raw := []byte(`{
		"daily_fees":"100",
		"daily_notional_volume":"1000000",
		"daily_premium_volume":"50000",
		"daily_trades":250,
		"open_interest":"500",
		"total_fees":"10000",
		"total_notional_volume":"100000000",
		"total_premium_volume":"500000",
		"total_trades":25000
	}`)
	var got types.Statistics
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "100", got.DailyFees.String())
	assert.Equal(t, "1000000", got.DailyNotionalVolume.String())
	assert.Equal(t, int64(250), got.DailyTrades)
	assert.Equal(t, "500", got.OpenInterest.String())
	assert.Equal(t, int64(25000), got.TotalTrades)
}
