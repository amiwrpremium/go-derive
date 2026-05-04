// Package contracts hosts on-chain helper interfaces — deposits,
// withdrawals, and session-key lifecycle — for Derive's smart-account
// model.
//
// # Status
//
// The package is intentionally a stub: the JSON-RPC layer
// ([github.com/amiwrpremium/go-derive/pkg/rest] and
// [github.com/amiwrpremium/go-derive/pkg/ws]) is sufficient to trade once
// collateral has been deposited via the Derive UI or another EVM tool.
// Every interface in this package is declared so that consumers can write
// code against the API today against a stable shape.
//
// All methods return [ErrNotImplemented].
package contracts
