package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestCandle_Decode(t *testing.T) {
	payload := `{
		"timestamp": 1700000000000,
		"open": "100",
		"high": "110",
		"low": "95",
		"close": "105",
		"volume": "10000"
	}`
	var c types.Candle
	require.NoError(t, json.Unmarshal([]byte(payload), &c))
	assert.Equal(t, "100", c.Open.String())
	assert.Equal(t, "110", c.High.String())
	assert.Equal(t, "95", c.Low.String())
	assert.Equal(t, "105", c.Close.String())
}

// Decimal's custom MarshalJSON defeats `omitempty` — a zero decimal still
// serialises as the string "0". Document and pin the actual behaviour.
func TestCandle_ZeroVolumeStillSerialized(t *testing.T) {
	in := types.Candle{
		Open:  types.MustDecimal("1"),
		High:  types.MustDecimal("1"),
		Low:   types.MustDecimal("1"),
		Close: types.MustDecimal("1"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"volume":"0"`)
}
