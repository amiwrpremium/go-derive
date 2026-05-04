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
