// Package netconf carries the network constants — endpoint URLs, chain IDs,
// and EIP-712 domain separators — for each Derive environment. These values
// are not user-tunable enums; they are concrete configuration that varies
// between mainnet and testnet.
package netconf

// Domain is the EIP-712 domain separator inputs needed by pkg/auth. Derive
// uses a per-network "ProtocolDomain" with name = "Matching" and version = "1.0"
// pinned to the matching engine contract.
type Domain struct {
	Name              string
	Version           string
	ChainID           int64
	VerifyingContract string
}

// EIP712Domain returns the domain separator inputs for the configured network.
//
// The Name/Version pair must exactly match the values the Matching.sol
// contract registers with its EIP712 base — currently "Matching" and "1.0".
// Drift here silently breaks signed-action verification with a 14014 error.
// pkg/auth/eip712_pinned_test.go pins the computed domain separator against
// the docs-published constant for both networks to catch any drift.
func (c Config) EIP712Domain() Domain {
	return Domain{
		Name:              "Matching",
		Version:           "1.0",
		ChainID:           c.ChainID,
		VerifyingContract: c.Contracts.MatchingEngine,
	}
}
