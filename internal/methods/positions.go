// Package methods — see collateral.go for the overview.
package methods

import (
	"context"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// GetPositions lists open positions on the subaccount. Private.
func (a *API) GetPositions(ctx context.Context) ([]types.Position, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	var resp struct {
		Positions []types.Position `json:"positions"`
	}
	err := a.call(ctx, "private/get_positions", params, &resp)
	return resp.Positions, err
}
