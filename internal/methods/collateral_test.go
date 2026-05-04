package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestGetCollateral_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_collaterals", map[string]any{
		"collaterals": []map[string]any{
			{"asset_name": "USDC", "asset_type": "erc20", "amount": "100", "mark_value": "100"},
		},
	})
	got, err := api.GetCollateral(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "USDC", got[0].AssetName)
	assert.Equal(t, "private/get_collaterals", ft.LastCall().Method)
}

func TestGetCollateral_Empty(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_collaterals", map[string]any{"collaterals": []any{}})
	got, err := api.GetCollateral(context.Background())
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestGetCollateral_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.GetCollateral(context.Background())
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}
