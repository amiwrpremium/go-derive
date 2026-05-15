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
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/internal/netconf"
)

// Signature is a 65-byte ECDSA signature in `r || s || v` byte order,
// where `v` follows Ethereum's 27/28 convention (not the raw 0/1 form
// go-ethereum produces internally — Derive's on-chain ecrecover path
// expects 27/28).
type Signature [65]byte

// Hex returns the canonical 0x-prefixed lowercase-hex representation.
// Length is always 132 characters (2 prefix + 65 bytes × 2).
func (s Signature) Hex() string {
	const hexChars = "0123456789abcdef"
	out := make([]byte, 2+len(s)*2)
	out[0] = '0'
	out[1] = 'x'
	for i, b := range s {
		out[2+i*2] = hexChars[b>>4]
		out[2+i*2+1] = hexChars[b&0x0f]
	}
	return string(out)
}

// Signer abstracts over the source of cryptographic signatures. The SDK
// uses it for both per-request auth-header signing (EIP-191) and
// per-action EIP-712 signing.
//
// Concrete implementations in this package:
//
//   - [LocalSigner]        — secp256k1 key held in process; owner == address.
//   - [SessionKeySigner]   — session key signs but reports a separate owner.
//
// External implementations are welcome: a hardware wallet, KMS-backed
// key, or HSM-backed key all fit cleanly behind this interface.
type Signer interface {
	// SessionAddress returns the public address whose signatures the
	// implementation produces. For session keys this is the session
	// key's own EOA; for [LocalSigner] it equals [Signer.OwnerAddress].
	SessionAddress() common.Address

	// OwnerAddress returns the owner (smart-account) address — the
	// long-lived wallet that authorised this signer. For [LocalSigner]
	// it equals [Signer.SessionAddress]; for [SessionKeySigner] it is
	// the distinct registered owner.
	OwnerAddress() common.Address

	// SignAction produces an EIP-712 signature over the action struct
	// hash with Derive's per-network domain. The implementation is
	// responsible for filling Action.Owner and Action.Signer if they
	// are zero.
	SignAction(ctx context.Context, domain netconf.Domain, action ActionData) (Signature, error)

	// SignAuthHeader produces an EIP-191 personal-sign signature over
	// the millisecond-timestamp string. The result is used as the
	// X-LyraSignature header on REST and as the `signature` field on
	// the WS `public/login` RPC.
	SignAuthHeader(ctx context.Context, ts time.Time) (Signature, error)
}
