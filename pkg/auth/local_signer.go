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
	"crypto/ecdsa"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"

	"github.com/amiwrpremium/go-derive/internal/netconf"
)

// LocalSigner holds an owner private key in process. It is the simplest
// Signer implementation; for production market-making prefer
// [SessionKeySigner] so the long-lived owner key never touches the process.
type LocalSigner struct {
	key  *ecdsa.PrivateKey
	addr common.Address
}

// NewLocalSigner parses a hex-encoded secp256k1 private key (with or without
// the 0x prefix) into a Signer.
func NewLocalSigner(hexKey string) (*LocalSigner, error) {
	hexKey = strings.TrimPrefix(hexKey, "0x")
	k, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, &derrors.SigningError{Op: "parse private key", Err: err}
	}
	return &LocalSigner{key: k, addr: crypto.PubkeyToAddress(k.PublicKey)}, nil
}

// Address returns the signer's public address.
func (s *LocalSigner) Address() common.Address { return s.addr }

// Owner returns the same address as Address — a LocalSigner has no separate owner.
func (s *LocalSigner) Owner() common.Address { return s.addr }

// SignAction signs Derive's EIP-712 Action struct.
func (s *LocalSigner) SignAction(_ context.Context, domain netconf.Domain, action ActionData) (Signature, error) {
	if action.Owner == (common.Address{}) {
		action.Owner = s.addr
	}
	if action.Signer == (common.Address{}) {
		action.Signer = s.addr
	}
	digest := hashTypedData(domain, action.Hash())
	return signDigest(s.key, digest)
}

// SignAuthHeader signs the millisecond timestamp string Derive expects.
func (s *LocalSigner) SignAuthHeader(_ context.Context, ts time.Time) (Signature, error) {
	msg := []byte(strconv.FormatInt(ts.UnixMilli(), 10))
	digest := hashEIP191(msg)
	return signDigest(s.key, digest)
}

// signDigest produces an Ethereum-flavoured 65-byte signature where the v
// byte is encoded as 27/28 (rather than 0/1) — the convention Derive expects
// in the on-chain ecrecover path.
func signDigest(key *ecdsa.PrivateKey, digest []byte) (Signature, error) {
	var sig Signature
	raw, err := crypto.Sign(digest, key)
	if err != nil {
		return sig, &derrors.SigningError{Op: "ecdsa sign", Err: err}
	}
	if len(raw) != 65 {
		return sig, &derrors.SigningError{Op: "ecdsa sign", Err: errShortSig}
	}
	copy(sig[:], raw)
	// go-ethereum returns v in {0, 1}; Solidity ecrecover wants {27, 28}.
	sig[64] += 27
	return sig, nil
}

var errShortSig = derrors.New("ecdsa signature too short")
