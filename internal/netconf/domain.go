// Package netconf carries the network constants — endpoint URLs, chain IDs,
// and EIP-712 domain separators — for each Derive environment. These values
// are not user-tunable enums; they are concrete configuration that varies
// between mainnet and testnet.
package netconf

// Domain is the EIP-712 domain separator inputs needed by pkg/auth. Derive
// uses a per-network "ProtocolDomain" with name = "Matching" and version = "1"
// pinned to the matching engine contract.
type Domain struct {
	Name              string
	Version           string
	ChainID           int64
	VerifyingContract string
}

// EIP712Domain returns the domain separator inputs for the configured network.
func (c Config) EIP712Domain() Domain {
	return Domain{
		Name:              "Matching",
		Version:           "1",
		ChainID:           c.ChainID,
		VerifyingContract: c.Contracts.MatchingEngine,
	}
}
