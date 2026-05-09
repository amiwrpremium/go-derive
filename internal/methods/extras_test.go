package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetPerpImpactTWAP_Raw covers the one method still surfaced as
// `json.RawMessage` by the SDK because it's not documented in the
// published OAS. The wrapper just forwards `params` and returns the
// raw payload.
func TestGetPerpImpactTWAP_Raw(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_perp_impact_twap", map[string]any{"impact_price": "100"})
	raw, err := api.GetPerpImpactTWAP(context.Background(), map[string]any{
		"currency": "BTC", "start_time": 0, "end_time": 1,
	})
	require.NoError(t, err)
	assert.JSONEq(t, `{"impact_price":"100"}`, string(raw))
}
