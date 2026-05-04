package codec_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/codec"
)

func TestPadLeft32_ShortInput(t *testing.T) {
	out := codec.PadLeft32([]byte{0x01, 0x02})
	assert.Len(t, out, 32)
	assert.Equal(t, byte(0x01), out[30])
	assert.Equal(t, byte(0x02), out[31])
	for i := 0; i < 30; i++ {
		assert.Equal(t, byte(0), out[i])
	}
}

func TestPadLeft32_ExactInput(t *testing.T) {
	in := bytes.Repeat([]byte{0xa5}, 32)
	out := codec.PadLeft32(in)
	assert.Equal(t, in, out)
}

func TestPadLeft32_EmptyInput(t *testing.T) {
	out := codec.PadLeft32(nil)
	assert.Len(t, out, 32)
	for _, b := range out {
		assert.Equal(t, byte(0), b)
	}
}

func TestPadLeft32_PanicsOnLong(t *testing.T) {
	assert.Panics(t, func() { codec.PadLeft32(make([]byte, 33)) })
}

func TestEncodeUint256_Zero(t *testing.T) {
	out, err := codec.EncodeUint256(big.NewInt(0))
	require.NoError(t, err)
	assert.True(t, bytes.Equal(out, make([]byte, 32)))
}

func TestEncodeUint256_Max(t *testing.T) {
	max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	out, err := codec.EncodeUint256(max)
	require.NoError(t, err)
	for _, b := range out {
		assert.Equal(t, byte(0xff), b)
	}
}

func TestEncodeUint256_Negative(t *testing.T) {
	_, err := codec.EncodeUint256(big.NewInt(-1))
	assert.Error(t, err)
}

func TestEncodeUint256_Overflow(t *testing.T) {
	_, err := codec.EncodeUint256(new(big.Int).Lsh(big.NewInt(1), 256))
	assert.Error(t, err)
}

func TestEncodeInt256_Zero(t *testing.T) {
	out, err := codec.EncodeInt256(big.NewInt(0))
	require.NoError(t, err)
	assert.True(t, bytes.Equal(out, make([]byte, 32)))
}

func TestEncodeInt256_PositiveSmall(t *testing.T) {
	out, err := codec.EncodeInt256(big.NewInt(7))
	require.NoError(t, err)
	assert.Equal(t, byte(7), out[31])
	for i := 0; i < 31; i++ {
		assert.Equal(t, byte(0), out[i])
	}
}

func TestEncodeInt256_NegativeOne_TwosComplement(t *testing.T) {
	out, err := codec.EncodeInt256(big.NewInt(-1))
	require.NoError(t, err)
	for _, b := range out {
		assert.Equal(t, byte(0xff), b)
	}
}

func TestEncodeInt256_LargePositive(t *testing.T) {
	// 2^200 fits in 255 bits.
	v := new(big.Int).Lsh(big.NewInt(1), 200)
	_, err := codec.EncodeInt256(v)
	assert.NoError(t, err)
}

func TestEncodeInt256_Overflow(t *testing.T) {
	v := new(big.Int).Lsh(big.NewInt(1), 256)
	_, err := codec.EncodeInt256(v)
	assert.Error(t, err)
}

func TestEncodeAddress_LeftPadded(t *testing.T) {
	a := common.HexToAddress("0x1111111111111111111111111111111111111111")
	out := codec.EncodeAddress(a)
	assert.Len(t, out, 32)
	for i := 0; i < 12; i++ {
		assert.Equal(t, byte(0), out[i])
	}
	assert.Equal(t, a.Bytes(), out[12:])
}

func TestEncodeAddress_ZeroAddress(t *testing.T) {
	out := codec.EncodeAddress(common.Address{})
	assert.Len(t, out, 32)
	for _, b := range out {
		assert.Equal(t, byte(0), b)
	}
}

func TestEncodeBytes32_Identity(t *testing.T) {
	in := bytes.Repeat([]byte{0x42}, 32)
	out := codec.EncodeBytes32(in)
	assert.Equal(t, in, out)
	assert.NotSame(t, &in, &out, "should return a copy")
}

func TestEncodeBytes32_PanicOnTooShort(t *testing.T) {
	assert.Panics(t, func() { codec.EncodeBytes32(make([]byte, 31)) })
}

func TestEncodeBytes32_PanicOnTooLong(t *testing.T) {
	assert.Panics(t, func() { codec.EncodeBytes32(make([]byte, 33)) })
}

func TestEncodeBytes32_EmptyInputPanics(t *testing.T) {
	assert.Panics(t, func() { codec.EncodeBytes32(nil) })
}
