package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
)

func TestGetPositions_Empty(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_positions", map[string]any{"positions": []any{}})
	got, err := api.GetPositions(context.Background())
	require.NoError(t, err)
	assert.Empty(t, got)
	assert.Equal(t, "private/get_positions", ft.LastCall().Method)
}

func TestGetPositions_NonEmpty(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_positions", map[string]any{
		"positions": []map[string]any{{
			"instrument_name": "BTC-PERP",
			"instrument_type": "perp",
			"amount":          "0.5",
			"average_price":   "65000",
			"mark_price":      "65500",
			"mark_value":      "32750",
			"unrealized_pnl":  "250",
			"realized_pnl":    "0",
		}},
	})
	got, err := api.GetPositions(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "BTC-PERP", got[0].InstrumentName)
}

func TestGetPositions_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.GetPositions(context.Background())
	assert.ErrorIs(t, err, derive.ErrSubaccountRequired)
}
