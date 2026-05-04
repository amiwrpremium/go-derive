package enums

// Environment selects which Derive deployment a client talks to. It is
// surfaced to users via the With* options on each client; the SDK turns
// it into the corresponding network configuration in
// [github.com/amiwrpremium/go-derive/internal/netconf].
type Environment string

const (
	// EnvironmentMainnet selects the production deployment (chain ID 957).
	EnvironmentMainnet Environment = "mainnet"
	// EnvironmentTestnet selects the demo/staging deployment (chain ID 901).
	EnvironmentTestnet Environment = "testnet"
)

// Valid reports whether the receiver is one of the defined environments.
func (e Environment) Valid() bool {
	switch e {
	case EnvironmentMainnet, EnvironmentTestnet:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (e Environment) Validate() error {
	if e.Valid() {
		return nil
	}
	return invalid("Environment", string(e))
}
