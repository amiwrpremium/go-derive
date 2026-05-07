// Package auth implements the two cryptographic signing flows Derive's
// API requires.
//
// # Two flows, one Signer
//
// Every authenticated Derive request involves cryptography in one of two
// places:
//
//  1. Per-request authentication of the caller. Sent as REST headers
//     (X-LyraWallet, X-LyraTimestamp, X-LyraSignature) or as a one-shot
//     `public/login` RPC over WebSocket. The signature is an EIP-191
//     personal-sign over the millisecond timestamp.
//
//  2. Per-action authorisation of order placement, cancels, transfers and
//     RFQ flows. The signature is an EIP-712 typed-data hash over an
//     `Action` struct whose `data` field is the keccak256 of an
//     ABI-encoded module-specific payload.
//
// Both flows go through the same [Signer] interface; concrete
// implementations include [LocalSigner] (owner key in process) and
// [SessionKeySigner] (session key delegating from a separate owner
// address).
//
// # Production setup
//
// Derive deployments use session keys. The owner is a smart-account on
// Derive Chain; the session key is a hot key registered on-chain as
// authorised to sign on its behalf. Use [NewSessionKeySigner] for
// production trading so the long-lived owner key never sits in the
// trading process's memory.
//
// # Test fixtures
//
// All signing test vectors live in pkg/auth/*_test.go. The tests verify
// that Derive's expected signature bytes can be reproduced from a known
// secp256k1 key — they're the canary for any future change to EIP-712
// hashing here or upstream in go-ethereum.
package auth

import (
	"sync"
	"sync/atomic"
	"time"
)

// NonceGen produces strictly-increasing nonces for action signing.
//
// Derive requires nonces to be unique per subaccount across an action's
// lifetime. This generator returns millisecond-timestamp-based nonces in
// the upper bits combined with a 16-bit incrementing suffix in the lower
// bits, which gives both human-readable ordering (the timestamp prefix)
// and collision resistance for many actions in the same millisecond.
//
// The zero value is not usable; construct via [NewNonceGen].
type NonceGen struct {
	last atomic.Uint64
	mu   sync.Mutex
	rand uint16
}

// NewNonceGen returns a generator seeded from the current time.
//
// The returned generator is safe for concurrent use.
func NewNonceGen() *NonceGen {
	g := &NonceGen{}
	// Take only the lower 16 bits of UnixNano() as the counter seed —
	// the explicit mask documents intent so gosec G115 sees the
	// narrowing as deliberate.
	g.rand = uint16(time.Now().UnixNano() & 0xFFFF)
	return g
}

// Next returns the next nonce.
//
// Under contention the algorithm bumps to (prev + 1) so the
// strict-monotonic property holds even when many goroutines call Next in
// the same millisecond.
func (g *NonceGen) Next() uint64 {
	g.mu.Lock()
	g.rand++
	suffix := uint64(g.rand)
	g.mu.Unlock()

	for {
		ms := uint64(time.Now().UnixMilli())
		candidate := ms<<16 | suffix
		prev := g.last.Load()
		if candidate <= prev {
			candidate = prev + 1
		}
		if g.last.CompareAndSwap(prev, candidate) {
			return candidate
		}
	}
}
