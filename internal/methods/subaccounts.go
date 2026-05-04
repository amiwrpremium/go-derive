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
