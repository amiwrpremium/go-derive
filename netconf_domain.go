// EIP-712 domain separator inputs live here. The [Domain] type is
// pulled out of netconf.go on its own — it is the bridge between
// [NetworkConfig] (network endpoints + contract addresses) and the
// signing path in auth.go, and lives separately for that reason.

package derive

// Domain is the EIP-712 domain separator inputs needed by the signing
// helpers in auth.go. Derive uses a per-network "ProtocolDomain" with
// name = "Matching" and version = "1" pinned to the matching engine
// contract.
type Domain struct {
	Name              string
	Version           string
	ChainID           int64
	VerifyingContract string
}

// EIP712Domain returns the domain separator inputs for the configured
// network.
func (c NetworkConfig) EIP712Domain() Domain {
	return Domain{
		Name:              "Matching",
		Version:           "1",
		ChainID:           c.ChainID,
		VerifyingContract: c.Contracts.MatchingEngine,
	}
}
