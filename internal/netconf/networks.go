// Package netconf carries the network constants — endpoint URLs, chain IDs,
// and EIP-712 domain separators — for each Derive environment. These values
// are not user-tunable enums; they are concrete configuration that varies
// between mainnet and testnet.
package netconf

import "fmt"

// Network identifies a Derive deployment.
type Network int

const (
	// NetworkUnknown is the zero value and is invalid.
	NetworkUnknown Network = iota
	// NetworkMainnet is Derive's production deployment.
	NetworkMainnet
	// NetworkTestnet is Derive's staging deployment.
	NetworkTestnet
)

// Config is the bundle of endpoints and chain parameters for one network.
type Config struct {
	Network     Network
	ChainID     int64
	HTTPURL     string
	WSURL       string
	ExplorerURL string
	// Contract addresses needed for EIP-712 signing.
	Contracts Contracts
}

// Contracts collects on-chain addresses referenced by the SDK. Only the
// matching engine and trade-module addresses are needed to compute action
// hashes for signing; deposit/withdraw addresses are used by pkg/contracts.
type Contracts struct {
	MatchingEngine string
	TradeModule    string
	DepositModule  string
	WithdrawModule string
	TransferModule string
}

// Mainnet returns Derive mainnet configuration.
//
// Endpoint addresses are the publicly documented values at
// https://docs.derive.xyz/. Update them here if Derive moves them.
func Mainnet() Config {
	return Config{
		Network:     NetworkMainnet,
		ChainID:     957,
		HTTPURL:     "https://api.lyra.finance",
		WSURL:       "wss://api.lyra.finance/ws",
		ExplorerURL: "https://explorer.lyra.finance",
		Contracts: Contracts{
			MatchingEngine: "0xB1dE3D5d4e1Fb9e60db9bf7F6F6F9b03F80cA0d8",
			TradeModule:    "0x87F2863866D85E3192a35A73b388BD625D83f2be",
			DepositModule:  "0x9B3FE5E5a3bcEa5df4E08c41Ce89C4e3Ff01Ace3",
			WithdrawModule: "0x9d0E8f5b25384C7310CB8C6aE32C8fbeb645d083",
			TransferModule: "0x01259207A40925b794C8ac320456F7F6c8FE2636",
		},
	}
}

// Testnet returns Derive testnet (staging) configuration.
func Testnet() Config {
	return Config{
		Network:     NetworkTestnet,
		ChainID:     901,
		HTTPURL:     "https://api-demo.lyra.finance",
		WSURL:       "wss://api-demo.lyra.finance/ws",
		ExplorerURL: "https://explorer-prod-testnet-0eakp60405.t.conduit.xyz",
		Contracts: Contracts{
			MatchingEngine: "0x6e1dF77Ade8Cd60F9F4F78a888F22Bd3aB52E0BC",
			TradeModule:    "0x87F2863866D85E3192a35A73b388BD625D83f2be",
			DepositModule:  "0x9B3FE5E5a3bcEa5df4E08c41Ce89C4e3Ff01Ace3",
			WithdrawModule: "0x9d0E8f5b25384C7310CB8C6aE32C8fbeb645d083",
			TransferModule: "0x01259207A40925b794C8ac320456F7F6c8FE2636",
		},
	}
}

// String implements fmt.Stringer for diagnostics.
func (n Network) String() string {
	switch n {
	case NetworkMainnet:
		return "mainnet"
	case NetworkTestnet:
		return "testnet"
	default:
		return fmt.Sprintf("unknown(%d)", int(n))
	}
}
