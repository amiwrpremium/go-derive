package auth_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func TestSignature_Hex_LengthAndPrefix(t *testing.T) {
	var s auth.Signature
	hex := s.Hex()
	assert.True(t, strings.HasPrefix(hex, "0x"))
	assert.Len(t, hex, 2+65*2)
}

func TestSignature_Hex_AllZeros(t *testing.T) {
	var s auth.Signature
	hex := s.Hex()
	assert.Equal(t, "0x"+strings.Repeat("00", 65), hex)
}

func TestSignature_Hex_AllOnes(t *testing.T) {
	var s auth.Signature
	for i := range s {
		s[i] = 0xff
	}
	assert.Equal(t, "0x"+strings.Repeat("ff", 65), s.Hex())
}

func TestSignature_Hex_MixedBytes(t *testing.T) {
	var s auth.Signature
	s[0] = 0x12
	s[1] = 0x34
	s[64] = 0xab
	hex := s.Hex()
	assert.True(t, strings.HasPrefix(hex, "0x1234"))
	assert.True(t, strings.HasSuffix(hex, "ab"))
}
