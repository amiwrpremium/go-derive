package auth_test

// Helpers shared by auth/* test files. The build tag is omitted so any test
// file in this package can call these.

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

// testKey is a throwaway secp256k1 key used only in tests.
const testKey = "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"

// personalHash mirrors auth/eip191.go for verification side.
func personalHash(msg []byte) []byte {
	prefix := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg)))
	h := crypto.NewKeccakState()
	_, _ = h.Write(prefix)
	_, _ = h.Write(msg)
	return h.Sum(nil)
}

// normaliseV converts Derive's v in {27,28} back to go-ethereum's {0,1}.
func normaliseV(sig []byte) []byte {
	out := make([]byte, len(sig))
	copy(out, sig)
	if len(out) == 65 {
		out[64] -= 27
	}
	return out
}
