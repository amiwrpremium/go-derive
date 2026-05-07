// Package auth.
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
