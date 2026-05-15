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

// SessionKeySigner wraps a [LocalSigner] (the session key) but reports the
// configured owner address as Owner(). This is the correct shape for Derive:
// orders are signed by the session key, but the smart account owner is the
// distinct on-chain wallet the session key was registered against.
type SessionKeySigner struct {
	inner *LocalSigner
	owner common.Address
}

// NewSessionKeySigner builds a SessionKeySigner from a hex session-key
// private key and the owner address it has been delegated by.
func NewSessionKeySigner(sessionHexKey string, owner common.Address) (*SessionKeySigner, error) {
	inner, err := NewLocalSigner(sessionHexKey)
	if err != nil {
		return nil, err
	}
	return &SessionKeySigner{inner: inner, owner: owner}, nil
}

// SessionAddress returns the session key's own EOA address.
func (s *SessionKeySigner) SessionAddress() common.Address { return s.inner.SessionAddress() }

// OwnerAddress returns the smart-account owner address.
func (s *SessionKeySigner) OwnerAddress() common.Address { return s.owner }

// SignAction populates Owner with the wallet owner address before signing.
func (s *SessionKeySigner) SignAction(ctx context.Context, domain netconf.Domain, action ActionData) (Signature, error) {
	action.Owner = s.owner
	action.Signer = s.inner.SessionAddress()
	return s.inner.SignAction(ctx, domain, action)
}

// SignAuthHeader signs as the session key.
func (s *SessionKeySigner) SignAuthHeader(ctx context.Context, ts time.Time) (Signature, error) {
	return s.inner.SignAuthHeader(ctx, ts)
}
