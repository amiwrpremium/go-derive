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
	"crypto/rand"
	"encoding/binary"
	"time"
)

// NonceGen produces nonces in the format Derive's engine documents:
//
//	nonce = (UTC milliseconds since epoch) * 1000 + random[0..999]
//
// That format keeps every nonce inside JSON's safe-integer range (2^53,
// the IEEE-754 double's max precise integer), so the engine's JSON
// parser doesn't truncate the value when it recomputes the EIP-712
// digest server-side. Any wider format silently breaks signature
// verification with a 14014 error.
//
// The shape and meaning match the official Python reference at
// derivexyz/v2-action-signing-python (utils.get_action_nonce) exactly.
//
// The engine only requires uniqueness, not strict monotonicity. The
// 1000-bucket random suffix gives ample collision headroom —
// birthday-paradox 50% probability at ~37 calls per millisecond, well
// above any non-HFT workload.
//
// The zero value is usable. NewNonceGen is kept for source
// compatibility with prior versions.
type NonceGen struct{}

// NewNonceGen returns a generator. The returned value is safe for
// concurrent use; crypto/rand.Read is goroutine-safe.
func NewNonceGen() *NonceGen { return &NonceGen{} }

// Next returns a fresh nonce in the documented format,
// unix_ms*1000 + suffix where suffix ∈ [0, 999]. The engine only
// requires uniqueness, not cryptographic secrecy — but the repo
// precedent (see internal/retry/backoff.go) is to source randomness
// from crypto/rand everywhere to satisfy the security linters
// without per-call nolint directives. The cost is one OS RNG read
// of 2 bytes per call, negligible relative to the network round-trip
// every signed action does anyway.
func (g *NonceGen) Next() uint64 {
	return uint64(time.Now().UTC().UnixMilli())*1000 + nonceSuffix()
}

// nonceSuffix returns a uniform value in [0, 999] sourced from
// crypto/rand. On the off chance crypto/rand.Read fails (it doesn't
// on Linux/macOS/Windows), this returns 0 — the resulting nonce is
// still valid because the ms-timestamp prefix dominates the
// collision space.
func nonceSuffix() uint64 {
	var b [2]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0
	}
	return uint64(binary.BigEndian.Uint16(b[:])) % 1000
}
