package auth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

// hashEIP191 computes the personal-sign digest:
//
//	keccak256( "\x19Ethereum Signed Message:\n" || len(msg) || msg )
//
// This is the ecrecover-compatible message hash used for Derive's REST and
// WebSocket auth signatures (the timestamp string is the message).
func hashEIP191(msg []byte) []byte {
	prefix := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg)))
	h := crypto.NewKeccakState()
	_, _ = h.Write(prefix)
	_, _ = h.Write(msg)
	return h.Sum(nil)
}
