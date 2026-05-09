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

// GetSubaccount fetches the configured subaccount snapshot. Private.
func (a *API) GetSubaccount(ctx context.Context) (derive.SubAccount, error) {
	if err := a.requireSubaccount(); err != nil {
		return derive.SubAccount{}, err
	}
	var sa derive.SubAccount
	err := a.call(ctx, "private/get_subaccount", map[string]any{
		"subaccount_id": a.Subaccount,
	}, &sa)
	return sa, err
}

// GetSubaccounts lists every subaccount owned by the wallet. Private.
func (a *API) GetSubaccounts(ctx context.Context) ([]int64, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	var resp struct {
		SubaccountIDs []int64 `json:"subaccount_ids"`
	}
	err := a.call(ctx, "private/get_subaccounts", map[string]any{
		"wallet": a.Signer.Owner().Hex(),
	}, &resp)
	return resp.SubaccountIDs, err
}
