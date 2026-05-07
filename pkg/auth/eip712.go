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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/amiwrpremium/go-derive/internal/codec"
	"github.com/amiwrpremium/go-derive/internal/netconf"
)

// keccak returns keccak256(b₀ || b₁ || ... || bₙ). It is a thin alias for
// readability at signing call sites.
func keccak(b ...[]byte) []byte {
	h := crypto.NewKeccakState()
	for _, p := range b {
		_, _ = h.Write(p)
	}
	return h.Sum(nil)
}

// domainSeparator computes the EIP-712 hashStruct(EIP712Domain) for a Derive
// network configuration.
//
// Derive uses the type
//
//	EIP712Domain(string name, string version, uint256 chainId,
//	             address verifyingContract)
//
// pinned to name="Matching", version="1", and the per-network matching
// engine address as the verifying contract.
func domainSeparator(d netconf.Domain) []byte {
	const typ = "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"
	typeHash := keccak([]byte(typ))
	nameHash := keccak([]byte(d.Name))
	versionHash := keccak([]byte(d.Version))

	chainID := bigInt(d.ChainID)
	chainIDBytes, _ := codec.EncodeUint256(chainID)
	verifying := codec.EncodeAddress(common.HexToAddress(d.VerifyingContract))

	return keccak(typeHash, nameHash, versionHash, chainIDBytes, verifying)
}

// hashTypedData applies the EIP-712 envelope:
//
//	keccak256( "\x19\x01" || domainSeparator || structHash )
//
// It is the digest signed by [Signer.SignAction].
func hashTypedData(domain netconf.Domain, structHash []byte) []byte {
	return keccak([]byte{0x19, 0x01}, domainSeparator(domain), structHash)
}
