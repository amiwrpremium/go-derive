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

// GetCollateral returns the collateral breakdown for the subaccount. Private.
func (a *API) GetCollateral(ctx context.Context) ([]types.Collateral, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var resp struct {
		Collaterals []types.Collateral `json:"collaterals"`
	}
	err := a.call(ctx, "private/get_collaterals", map[string]any{
		"subaccount_id": a.Subaccount,
	}, &resp)
	return resp.Collaterals, err
}
