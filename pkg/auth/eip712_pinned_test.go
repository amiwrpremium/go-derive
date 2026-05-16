package auth

// Pin the EIP-712 domain separator and action typehash to the values
// published in docs.derive.xyz/reference/protocol-constants. These tests
// caught three signing bugs that survived the existing
// sensitivity / determinism tests:
//
//   - domain version "1" vs the contract's "1.0"
//   - mainnet/testnet MatchingEngine (= verifyingContract) addresses
//   - actionTypeHash declaring `bytes32 data` vs the contract's `bytes data`
//
// If any of these tests fail after a future change, the signed action
// digest no longer matches what the engine recomputes server-side and
// every signed action gets rejected with code 14014 ("Signature invalid
// for message or transaction"). Reproduce the doc values by reading
// `docs.derive.xyz/reference/protocol-constants` (HTML or .md).
//
// Lives in package auth so it can read the private actionTypeHash var
// and call the private domainSeparator function.

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/internal/netconf"
)

func TestActionTypeHash_MatchesDocPublished(t *testing.T) {
	const want = "4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17"
	got := hex.EncodeToString(actionTypeHash)
	assert.Equal(t, want, got,
		"actionTypeHash drift: the Action type string in action.go no longer "+
			"hashes to the value the Matching.sol contract uses. Every signed "+
			"action will be rejected by the engine with code 14014.")
}

func TestDomainSeparator_Mainnet_MatchesDocPublished(t *testing.T) {
	const want = "d96e5f90797da7ec8dc4e276260c7f3f87fedf68775fbe1ef116e996fc60441b"
	got := hex.EncodeToString(domainSeparator(netconf.Mainnet().EIP712Domain()))
	assert.Equal(t, want, got,
		"mainnet domain separator drift: one of (Name, Version, ChainID, "+
			"MatchingEngine) no longer matches what Matching.sol on mainnet "+
			"registered with its EIP-712 base. Every signed action on mainnet "+
			"will be rejected with code 14014.")
}

func TestDomainSeparator_Testnet_MatchesDocPublished(t *testing.T) {
	const want = "9bcf4dc06df5d8bf23af818d5716491b995020f377d3b7b64c29ed14e3dd1105"
	got := hex.EncodeToString(domainSeparator(netconf.Testnet().EIP712Domain()))
	assert.Equal(t, want, got,
		"testnet domain separator drift — same class of bug as mainnet "+
			"variant above. Every signed action on testnet will be rejected "+
			"with code 14014.")
}
