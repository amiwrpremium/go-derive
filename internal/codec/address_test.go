package codec_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/codec"
)

func TestParseAddress_HexLower(t *testing.T) {
	got, err := codec.ParseAddress("0x1111111111111111111111111111111111111111")
	require.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0x1111111111111111111111111111111111111111"), got)
}

func TestParseAddress_HexUpper(t *testing.T) {
	got, err := codec.ParseAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	require.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), got)
}

func TestParseAddress_RejectsInvalid(t *testing.T) {
	for _, in := range []string{"", "not-an-address", "0xabc", "12345"} {
		t.Run(in, func(t *testing.T) {
			_, err := codec.ParseAddress(in)
			assert.Error(t, err)
		})
	}
}
