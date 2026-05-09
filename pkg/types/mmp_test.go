package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestMMPConfigResult_Decode(t *testing.T) {
	raw := []byte(`{
		"subaccount_id":42,
		"currency":"BTC",
		"mmp_frozen_time":5000,
		"mmp_interval":1000,
		"mmp_amount_limit":"100",
		"mmp_delta_limit":"50",
		"mmp_unfreeze_time":1700000060000,
		"is_frozen":false
	}`)
	var got types.MMPConfigResult
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, int64(42), got.SubaccountID)
	assert.Equal(t, "BTC", got.Currency)
	assert.Equal(t, int64(5000), got.MMPFrozenTime)
	assert.Equal(t, int64(1000), got.MMPInterval)
	assert.Equal(t, "100", got.MMPAmountLimit.String())
	assert.Equal(t, "50", got.MMPDeltaLimit.String())
	assert.Equal(t, int64(1700000060000), got.MMPUnfreezeTime)
	assert.False(t, got.IsFrozen)
}

func TestMMPConfigResult_OptionalAmountLimitsOmitted(t *testing.T) {
	raw := []byte(`{
		"subaccount_id":1,
		"currency":"ETH",
		"mmp_frozen_time":0,
		"mmp_interval":0,
		"mmp_unfreeze_time":0,
		"is_frozen":false
	}`)
	var got types.MMPConfigResult
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "0", got.MMPAmountLimit.String())
	assert.Equal(t, "0", got.MMPDeltaLimit.String())
}
