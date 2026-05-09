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

import (
	"context"
	"time"

	"github.com/amiwrpremium/go-derive"
)

// SessionKeyManager is the contract for the session-key lifecycle.
// Session keys are addresses authorised to sign Derive actions on behalf
// of the owner wallet; they limit blast radius if a hot key is
// compromised.
type SessionKeyManager interface {
	// Register adds a session key authorised to sign actions on behalf of
	// the owner wallet, valid until expiry. It returns the registration
	// transaction hash on success.
	//
	// Returns [ErrNotImplemented].
	Register(ctx context.Context, sessionKey derive.Address, expiry time.Time) (derive.TxHash, error)

	// Revoke immediately deauthorises a session key. It returns the
	// revocation transaction hash on success.
	//
	// Returns [ErrNotImplemented].
	Revoke(ctx context.Context, sessionKey derive.Address) (derive.TxHash, error)
}
