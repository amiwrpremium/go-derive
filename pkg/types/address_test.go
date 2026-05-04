package types_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

const (
	addrUpper = "0x1234567890ABCDEF1234567890ABCDEF12345678"
	addrLower = "0x1234567890abcdef1234567890abcdef12345678"
)

func TestNewAddress_Valid(t *testing.T) {
	a, err := types.NewAddress(addrLower)
	require.NoError(t, err)
	assert.False(t, a.IsZero())
	assert.Equal(t, common.HexToAddress(addrLower), a.Common())
}

func TestNewAddress_Empty(t *testing.T) {
	a, err := types.NewAddress("")
	require.NoError(t, err)
	assert.True(t, a.IsZero())
}

func TestNewAddress_Invalid(t *testing.T) {
	_, err := types.NewAddress("not-an-address")
	assert.Error(t, err)
}

func TestMustAddress_Panics(t *testing.T) {
	assert.Panics(t, func() { types.MustAddress("nope") })
	assert.NotPanics(t, func() { types.MustAddress(addrLower) })
}

func TestAddress_StringIsChecksummed(t *testing.T) {
	a := types.MustAddress(addrLower)
	// EIP-55 mixed-case
	assert.Equal(t, common.HexToAddress(addrLower).Hex(), a.String())
}

func TestAddress_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		A types.Address `json:"a"`
	}
	in := wrap{A: types.MustAddress(addrLower)}
	b, err := json.Marshal(in)
	require.NoError(t, err)

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.A, out.A)
}

func TestAddress_UnmarshalEmpty(t *testing.T) {
	var a types.Address
	require.NoError(t, json.Unmarshal([]byte(`""`), &a))
	assert.True(t, a.IsZero())
}

func TestAddress_UnmarshalInvalid(t *testing.T) {
	var a types.Address
	err := json.Unmarshal([]byte(`"not-an-address"`), &a)
	assert.Error(t, err)
}

func TestAddress_UnmarshalNonString(t *testing.T) {
	var a types.Address
	err := json.Unmarshal([]byte(`123`), &a)
	assert.Error(t, err)
}
