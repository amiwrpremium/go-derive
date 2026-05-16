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
	"math/rand/v2"
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
// The engine only requires uniqueness, not strict monotonicity (the
// Python reference uses non-monotonic `random.randint(0, 999)` and
// works fine in production). The 1000-bucket suffix gives ample
// collision headroom — birthday-paradox 50% probability at ~37
// calls per millisecond, well above any non-HFT workload.
//
// The zero value is usable. NewNonceGen is kept for source
// compatibility with prior versions.
type NonceGen struct{}

// NewNonceGen returns a generator. The returned value is safe for
// concurrent use; math/rand/v2's default Source has internal locking
// and is auto-seeded from crypto/rand at package init.
func NewNonceGen() *NonceGen { return &NonceGen{} }

// Next returns a fresh nonce in the documented format,
// unix_ms*1000 + suffix where suffix ∈ [0, 999]. The engine only
// requires uniqueness, not cryptographic secrecy, so math/rand/v2
// is the right tool — using crypto/rand here would buy nothing and
// cost a syscall per call.
func (g *NonceGen) Next() uint64 {
	//nolint:gosec // G404: see doc comment — nonce uniqueness only.
	return uint64(time.Now().UTC().UnixMilli())*1000 + uint64(rand.IntN(1000))
}
