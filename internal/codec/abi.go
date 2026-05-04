package codec

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// PadLeft32 left-pads b to 32 bytes. Used for ABI-style encoding of action
// module data. It panics if b is longer than 32 bytes since that indicates a
// caller bug.
func PadLeft32(b []byte) []byte {
	if len(b) > 32 {
		panic(fmt.Sprintf("codec: PadLeft32 input too long (%d)", len(b)))
	}
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}

// EncodeUint256 ABI-encodes a uint256 as 32 bytes big-endian. It returns an
// error if n is negative or larger than 256 bits.
func EncodeUint256(n *big.Int) ([]byte, error) {
	if n.Sign() < 0 {
		return nil, fmt.Errorf("codec: EncodeUint256: negative value")
	}
	if n.BitLen() > 256 {
		return nil, fmt.Errorf("codec: EncodeUint256: overflow (%d bits)", n.BitLen())
	}
	return PadLeft32(n.Bytes()), nil
}

// EncodeInt256 ABI-encodes an int256 as 32 bytes two's-complement big-endian.
func EncodeInt256(n *big.Int) ([]byte, error) {
	if n.BitLen() > 255 && n.Sign() != 0 {
		return nil, fmt.Errorf("codec: EncodeInt256: overflow (%d bits)", n.BitLen())
	}
	if n.Sign() >= 0 {
		return PadLeft32(n.Bytes()), nil
	}
	// Two's complement for negative numbers: add 2^256.
	mod := new(big.Int).Lsh(big.NewInt(1), 256)
	twoC := new(big.Int).Add(n, mod)
	return PadLeft32(twoC.Bytes()), nil
}

// EncodeAddress ABI-encodes a 20-byte address as 32 bytes left-padded.
func EncodeAddress(a common.Address) []byte {
	return PadLeft32(a.Bytes())
}

// EncodeBytes32 returns b as-is if it is exactly 32 bytes, otherwise panics —
// callers should always pass keccak256 hashes.
func EncodeBytes32(b []byte) []byte {
	if len(b) != 32 {
		panic(fmt.Sprintf("codec: EncodeBytes32 expects 32 bytes, got %d", len(b)))
	}
	out := make([]byte, 32)
	copy(out, b)
	return out
}
