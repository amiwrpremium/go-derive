// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the bridge-balance shape returned by
// `public/get_bridge_balances`.
package types

// BridgeBalance is one entry in `public/get_bridge_balances` —
// describes the on-chain balance held by Derive's cross-chain bridge
// for one (chain, integrator) tuple.
type BridgeBalance struct {
	// Name is the bridge / integration's display name.
	Name string `json:"name"`
	// Integrator is the integrator name (e.g. the bridge provider).
	Integrator string `json:"integrator,omitempty"`
	// ChainID is the chain the balance is held on.
	ChainID int64 `json:"chain_id"`
	// Balance is the current $ balance on the bridge.
	Balance Decimal `json:"balance"`
	// BalanceHours is the projected hours of bridge runway given
	// current outflow rate.
	BalanceHours Decimal `json:"balance_hours,omitempty"`
}
