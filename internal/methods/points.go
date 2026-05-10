// Package methods is the shared implementation of every JSON-RPC method
// Derive exposes. Both pkg/rest.Client and pkg/ws.Client embed *API so that
// each method is defined exactly once, parameterised by the underlying
// transport.
//
// Public methods are unauthenticated; private methods require Signer to be
// non-nil. Private methods that mutate orders also use the Domain to sign
// the per-action EIP-712 hash.
package methods

import (
	"context"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// GetAllPoints returns the program-wide points snapshot for one
// program: aggregate notional volume, user count, and per-wallet
// points map. Public.
//
// The per-wallet `points` map is preserved as raw JSON because the
// inner schema varies per program; decode further at the call site.
func (a *API) GetAllPoints(ctx context.Context, program string) (*types.AllPointsResult, error) {
	var resp types.AllPointsResult
	if err := a.call(ctx, "public/get_all_points", map[string]any{"program": program}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
