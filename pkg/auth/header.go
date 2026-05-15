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
	"net/http"
	"strconv"
	"time"
)

// HTTPHeaders builds the per-request authentication headers Derive expects
// on every REST call:
//
//	X-LyraWallet     — the owner address as 0x-prefixed hex
//	X-LyraTimestamp  — the current time as milliseconds since the Unix epoch
//	X-LyraSignature  — the EIP-191 signature over the timestamp string
//
// Despite the rename to "Derive", the header names retain their "Lyra"
// prefix server-side.
//
// If signer is nil, HTTPHeaders returns (nil, nil) — used by the public-only
// path of the HTTP transport. Errors from [Signer.SignAuthHeader] are
// propagated unmodified.
func HTTPHeaders(ctx context.Context, signer Signer, now time.Time) (http.Header, error) {
	if signer == nil {
		return nil, nil
	}
	sig, err := signer.SignAuthHeader(ctx, now)
	if err != nil {
		return nil, err
	}
	h := make(http.Header, 3)
	h.Set("X-LyraWallet", signer.OwnerAddress().Hex())
	h.Set("X-LyraTimestamp", strconv.FormatInt(now.UnixMilli(), 10))
	h.Set("X-LyraSignature", sig.Hex())
	return h, nil
}
