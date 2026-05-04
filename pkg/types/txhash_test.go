package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

const sampleHash = "0x1111111111111111111111111111111111111111111111111111111111111111"

func TestNewTxHash_Valid(t *testing.T) {
	h, err := types.NewTxHash(sampleHash)
	require.NoError(t, err)
	assert.False(t, h.IsZero())
	assert.Equal(t, sampleHash, h.String())
}

func TestNewTxHash_Empty(t *testing.T) {
	h, err := types.NewTxHash("")
	require.NoError(t, err)
	assert.True(t, h.IsZero())
}

func TestNewTxHash_Invalid(t *testing.T) {
	_, err := types.NewTxHash("0xabc")
	assert.Error(t, err)
	_, err = types.NewTxHash("notahash")
	assert.Error(t, err)
}

func TestTxHash_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		H types.TxHash `json:"h"`
	}
	in := wrap{H: must(types.NewTxHash(sampleHash))}
	b, err := json.Marshal(in)
	require.NoError(t, err)

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.H, out.H)
}

func TestTxHash_UnmarshalEmpty(t *testing.T) {
	var h types.TxHash
	require.NoError(t, json.Unmarshal([]byte(`""`), &h))
	assert.True(t, h.IsZero())
}

func TestTxHash_UnmarshalInvalid(t *testing.T) {
	var h types.TxHash
	assert.Error(t, json.Unmarshal([]byte(`"0xabc"`), &h))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
