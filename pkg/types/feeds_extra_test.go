package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestInterestRateHistoryItem_Decode(t *testing.T) {
	raw := []byte(`{
		"block": 12345,
		"timestamp_sec": 1700000000,
		"borrow_apy": "0.08",
		"supply_apy": "0.04",
		"total_borrow": "5000000",
		"total_supply": "10000000"
	}`)
	var i types.InterestRateHistoryItem
	require.NoError(t, json.Unmarshal(raw, &i))
	assert.Equal(t, int64(12345), i.Block)
	assert.Equal(t, int64(1700000000), i.TimestampSec)
	assert.Equal(t, "0.08", i.BorrowAPY.String())
	assert.Equal(t, "10000000", i.TotalSupply.String())
}
