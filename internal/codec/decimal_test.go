package codec_test

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/codec"
)

func TestDecimalToU256_PositiveExact(t *testing.T) {
	d := decimal.RequireFromString("1.5")
	got, err := codec.DecimalToU256(d)
	require.NoError(t, err)
	want, _ := new(big.Int).SetString("1500000000000000000", 10)
	assert.Equal(t, 0, want.Cmp(got))
}

func TestDecimalToU256_Zero(t *testing.T) {
	got, err := codec.DecimalToU256(decimal.Zero)
	require.NoError(t, err)
	assert.Equal(t, 0, big.NewInt(0).Cmp(got))
}

func TestDecimalToU256_NegativeRejected(t *testing.T) {
	_, err := codec.DecimalToU256(decimal.RequireFromString("-1"))
	assert.Error(t, err)
}

func TestDecimalToU256_PrecisionExceeded(t *testing.T) {
	_, err := codec.DecimalToU256(decimal.RequireFromString("0.0000000000000000005"))
	assert.Error(t, err)
}

func TestDecimalToU256_MaxRepresentable(t *testing.T) {
	// Just under 2^192 (well under 2^256) — should succeed.
	d := decimal.RequireFromString("100000000000000000000")
	_, err := codec.DecimalToU256(d)
	assert.NoError(t, err)
}

func TestDecimalToI256_NegativeAccepted(t *testing.T) {
	got, err := codec.DecimalToI256(decimal.RequireFromString("-2.5"))
	require.NoError(t, err)
	want, _ := new(big.Int).SetString("-2500000000000000000", 10)
	assert.Equal(t, 0, want.Cmp(got))
}

func TestDecimalToI256_Positive(t *testing.T) {
	got, err := codec.DecimalToI256(decimal.RequireFromString("3"))
	require.NoError(t, err)
	want := new(big.Int)
	want.SetString("3000000000000000000", 10)
	assert.Equal(t, 0, want.Cmp(got))
}

func TestDecimalToI256_Zero(t *testing.T) {
	got, err := codec.DecimalToI256(decimal.Zero)
	require.NoError(t, err)
	assert.Equal(t, 0, big.NewInt(0).Cmp(got))
}

func TestDecimalToI256_PrecisionExceeded(t *testing.T) {
	_, err := codec.DecimalToI256(decimal.RequireFromString("0.0000000000000000005"))
	assert.Error(t, err)
}
