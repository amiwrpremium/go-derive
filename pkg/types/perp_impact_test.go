package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestPerpImpactTWAP_Decode(t *testing.T) {
	raw := []byte(`{
		"currency":"BTC",
		"mid_price_diff_twap":"0.5",
		"ask_impact_diff_twap":"1.2",
		"bid_impact_diff_twap":"-0.8"
	}`)
	var got types.PerpImpactTWAP
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "BTC", got.Currency)
	assert.Equal(t, "0.5", got.MidPriceDiffTWAP.String())
	assert.Equal(t, "1.2", got.AskImpactDiffTWAP.String())
	assert.Equal(t, "-0.8", got.BidImpactDiffTWAP.String())
}

func TestPerpImpactTWAP_QuietBook(t *testing.T) {
	// In a quiet book the engine reports all three diffs as zero.
	raw := []byte(`{"currency":"BTC","mid_price_diff_twap":"0","ask_impact_diff_twap":"0","bid_impact_diff_twap":"0"}`)
	var got types.PerpImpactTWAP
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "0", got.MidPriceDiffTWAP.String())
}
