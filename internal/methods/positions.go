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

	"github.com/amiwrpremium/go-derive"
)

// GetPositions lists open positions on the subaccount. Private.
func (a *API) GetPositions(ctx context.Context) ([]derive.Position, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	var resp struct {
		Positions []derive.Position `json:"positions"`
	}
	err := a.call(ctx, "private/get_positions", params, &resp)
	return resp.Positions, err
}
