// Package auth.
package auth

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/internal/netconf"
)

// SessionKeySigner wraps a [LocalSigner] (the session key) but reports the
// configured owner address as Owner(). This is the correct shape for Derive:
// orders are signed by the session key, but the smart account owner is the
// distinct on-chain wallet the session key was registered against.
type SessionKeySigner struct {
	inner *LocalSigner
	owner common.Address
}

// NewSessionKeySigner builds a SessionKeySigner from a hex session-key
// private key and the owner address it has been delegated by.
func NewSessionKeySigner(sessionHexKey string, owner common.Address) (*SessionKeySigner, error) {
	inner, err := NewLocalSigner(sessionHexKey)
	if err != nil {
		return nil, err
	}
	return &SessionKeySigner{inner: inner, owner: owner}, nil
}

// Address returns the session key address.
func (s *SessionKeySigner) Address() common.Address { return s.inner.Address() }

// Owner returns the smart account owner address.
func (s *SessionKeySigner) Owner() common.Address { return s.owner }

// SignAction populates Owner with the wallet owner address before signing.
func (s *SessionKeySigner) SignAction(ctx context.Context, domain netconf.Domain, action ActionData) (Signature, error) {
	action.Owner = s.owner
	action.Signer = s.inner.Address()
	return s.inner.SignAction(ctx, domain, action)
}

// SignAuthHeader signs as the session key.
func (s *SessionKeySigner) SignAuthHeader(ctx context.Context, ts time.Time) (Signature, error) {
	return s.inner.SignAuthHeader(ctx, ts)
}
