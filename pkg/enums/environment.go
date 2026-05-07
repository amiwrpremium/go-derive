// Package enums declares the named-string enums used across the SDK.
//
// Each enum is a defined string type — the simplest idiom in Go that gives
// you exhaustive switch warnings, free JSON round-trips, and a domain-specific
// receiver set without the heavyweight ceremony of an `iota` block plus
// custom marshalers. Aliases of underlying string types like:
//
//	type Direction string
//	const DirectionBuy Direction = "buy"
//
// match what big Go SDKs (aws-sdk-go-v2, stripe-go) use, and the wire format
// they produce is byte-for-byte what Derive expects.
//
// Every enum exposes a Valid method for cheap input validation. Some, like
// [Direction], expose extra domain helpers ([Direction.Sign],
// [Direction.Opposite], [OrderStatus.Terminal]).
//
// # Validating untrusted input
//
// Always check [Direction.Valid] (or the corresponding Valid method on the
// enum) before passing user-provided strings into the SDK. The Go type
// system can't prevent constructing an out-of-range value via `Direction("x")`,
// so the runtime check is the safety net.
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
