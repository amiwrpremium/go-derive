// Package contracts.
package contracts

import (
	"context"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/types"
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
	Register(ctx context.Context, sessionKey types.Address, expiry time.Time) (types.TxHash, error)

	// Revoke immediately deauthorises a session key. It returns the
	// revocation transaction hash on success.
	//
	// Returns [ErrNotImplemented].
	Revoke(ctx context.Context, sessionKey types.Address) (types.TxHash, error)
}
