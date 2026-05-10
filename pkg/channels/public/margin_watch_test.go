package public_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/channels/public"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestMarginWatch_Name(t *testing.T) {
	assert.Equal(t, "margin_watch", public.MarginWatch{}.Name())
}

func TestMarginWatch_Decode(t *testing.T) {
	raw := json.RawMessage(`[
		{
			"subaccount_id": 42,
			"currency": "USDC",
			"margin_type": "PM",
			"subaccount_value": "10000",
			"maintenance_margin": "-5",
			"valuation_timestamp": 1700000000
		},
		{
			"subaccount_id": 7,
			"currency": "USDC",
			"margin_type": "SM",
			"subaccount_value": "100",
			"maintenance_margin": "10",
			"valuation_timestamp": 1700000000
		}
	]`)
	got, err := public.MarginWatch{}.Decode(raw)
	require.NoError(t, err)
	events, ok := got.([]types.MarginWatch)
	require.True(t, ok, "Decode must return []types.MarginWatch")
	require.Len(t, events, 2)
	assert.Equal(t, int64(42), events[0].SubaccountID)
	assert.Equal(t, enums.MarginTypePM, events[0].MarginType)
	assert.Equal(t, "-5", events[0].MaintenanceMargin.String(),
		"negative maintenance margin signals at-risk")
	assert.Equal(t, enums.MarginTypeSM, events[1].MarginType)
}
