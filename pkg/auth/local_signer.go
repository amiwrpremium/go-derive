// Package auth.
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
