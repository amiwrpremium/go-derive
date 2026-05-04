package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigInt(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0"},
		{1, "1"},
		{-1, "-1"},
		{1<<62 - 1, "4611686018427387903"},
		{-(1 << 62), "-4611686018427387904"},
	}
	for _, c := range cases {
		got := bigInt(c.in)
		assert.Equal(t, c.want, got.String(), "bigInt(%d)", c.in)
	}
}

func TestBigInt_DistinctPointers(t *testing.T) {
	a := bigInt(1)
	b := bigInt(1)
	a.SetInt64(2)
	assert.Equal(t, "1", b.String(), "bigInt must return independent *big.Int values")
	// belt-and-suspenders: confirm a moved.
	assert.Equal(t, "2", a.String())
}

func TestBigUint(t *testing.T) {
	// bigUint must NOT panic on values that overflow int64. The largest
	// uint64 (2^64 - 1) is the canonical edge case.
	got := bigUint(^uint64(0))
	assert.Equal(t, "18446744073709551615", got.String())
}
