package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestMakerProgram_Decode(t *testing.T) {
	raw := []byte(`{
		"name": "options-q1",
		"asset_types": ["option"],
		"currencies": ["ETH","BTC"],
		"min_notional": "10000",
		"rewards": {"DRV":"100000","OP":"5000"},
		"start_timestamp": 1700000000,
		"end_timestamp": 1702592000
	}`)
	var p types.MakerProgram
	require.NoError(t, json.Unmarshal(raw, &p))
	assert.Equal(t, "options-q1", p.Name)
	assert.Equal(t, []string{"option"}, p.AssetTypes)
	assert.Equal(t, "10000", p.MinNotional.String())
	assert.Equal(t, int64(1700000000), p.StartTimestamp)
	assert.Equal(t, "100000", p.Rewards["DRV"].String())
}

func TestMakerProgramScore_Decode(t *testing.T) {
	raw := []byte(`{
		"program": {
			"name":"options-q1","asset_types":["option"],"currencies":["ETH"],
			"min_notional":"10000","rewards":{"DRV":"100"},
			"start_timestamp":1700000000,"end_timestamp":1702592000
		},
		"scores": [{
			"wallet":"0x1111111111111111111111111111111111111111",
			"coverage_score":"0.8","quality_score":"0.9","holder_boost":"1.1",
			"volume":"50000","volume_multiplier":"1","total_score":"0.792"
		}],
		"total_score": "1.5",
		"total_volume": "100000"
	}`)
	var s types.MakerProgramScore
	require.NoError(t, json.Unmarshal(raw, &s))
	assert.Equal(t, "options-q1", s.Program.Name)
	assert.Equal(t, "1.5", s.TotalScore.String())
	require.Len(t, s.Scores, 1)
	assert.Equal(t, "0x1111111111111111111111111111111111111111", s.Scores[0].Wallet.String())
	assert.Equal(t, "0.792", s.Scores[0].TotalScore.String())
}

func TestReferralPerformance_Decode(t *testing.T) {
	raw := []byte(`{
		"referral_code": "ALICE",
		"fee_share_percentage": "0.2",
		"stdrv_balance": "1000",
		"total_notional_volume": "500000",
		"total_referred_fees": "1500",
		"total_fee_rewards": "300",
		"total_builder_fee_collected": "50",
		"rewards": {
			"taker": {"ETH": {"perp": {
				"notional_volume":"500000","referred_fee":"1500","fee_reward":"300",
				"builder_fee":"50","unique_traders_referred":7
			}}}
		}
	}`)
	var r types.ReferralPerformance
	require.NoError(t, json.Unmarshal(raw, &r))
	assert.Equal(t, "ALICE", r.ReferralCode)
	assert.Equal(t, "0.2", r.FeeSharePercentage.String())
	leaf := r.Rewards["taker"]["ETH"]["perp"]
	assert.Equal(t, "500000", leaf.NotionalVolume.String())
	assert.Equal(t, int64(7), leaf.UniqueTradersReferred)
}
