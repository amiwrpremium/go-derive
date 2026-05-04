package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestLiquidation_Decode(t *testing.T) {
	payload := `{
		"subaccount_id": 9,
		"timestamp": 1700000000000,
		"tx_hash": "0x1111111111111111111111111111111111111111111111111111111111111111"
	}`
	var l types.Liquidation
	require.NoError(t, json.Unmarshal([]byte(payload), &l))
	assert.Equal(t, int64(9), l.SubaccountID)
	assert.False(t, l.TxHash.IsZero())
}

// Custom MarshalJSON on TxHash defeats omitempty (Go's json package only
// omits the natural zero value for built-in types). The wire format always
// includes the hash field even when zero — document and pin that.
func TestLiquidation_ZeroTxHashSerializedAsZeroString(t *testing.T) {
	in := types.Liquidation{SubaccountID: 1}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"tx_hash":"0x0000000000000000000000000000000000000000000000000000000000000000"`)
}
