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

// GetSubaccount fetches the configured subaccount snapshot. Private.
func (a *API) GetSubaccount(ctx context.Context) (types.SubAccount, error) {
	if err := a.requireSubaccount(); err != nil {
		return types.SubAccount{}, err
	}
	var sa types.SubAccount
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

// ChangeSubaccountLabel sets the human-readable label on the
// configured subaccount. Private.
//
// The endpoint returns no useful data on success; this method
// surfaces a `nil` error.
func (a *API) ChangeSubaccountLabel(ctx context.Context, label string) error {
	if err := a.requireSigner(); err != nil {
		return err
	}
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	return a.call(ctx, "private/change_subaccount_label", map[string]any{
		"subaccount_id": a.Subaccount,
		"label":         label,
	}, nil)
}
