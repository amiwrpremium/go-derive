package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestFundingRateHistoryItem_Decode(t *testing.T) {
	raw := []byte(`{"timestamp":1700000000000,"funding_rate":"0.000125"}`)
	var got types.FundingRateHistoryItem
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, int64(1700000000000), got.Timestamp.Millis())
	assert.Equal(t, "0.000125", got.FundingRate.String())
}

func TestSpotFeedHistoryItem_Decode(t *testing.T) {
	raw := []byte(`{
		"timestamp":1700000000000,
		"timestamp_bucket":1700000000000,
		"price":"50000.5"
	}`)
	var got types.SpotFeedHistoryItem
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "50000.5", got.Price.String())
}

func TestSignedFeeds_Decode(t *testing.T) {
	raw := []byte(`{
		"spot_data":{
			"BTC":{"currency":"BTC","price":"50000","confidence":"0.99","timestamp":1700000000000,"deadline":1700000060000,"feed_source_type":"S","signatures":{"signers":["0xa"],"signatures":["0xs"]}}
		},
		"perp_data":{
			"BTC":{"P":{"currency":"BTC","type":"P","spot_diff_value":"100","confidence":"0.99","timestamp":1700000000000,"deadline":1700000060000,"signatures":{"signers":["0xa"],"signatures":["0xs"]}}}
		},
		"fwd_data":{
			"BTC":{"1700000000":{"currency":"BTC","expiry":1700000000,"fwd_diff":"50","spot_aggregate_latest":"50050","spot_aggregate_start":"50000","confidence":"0.95","timestamp":1700000000000,"deadline":1700000060000,"signatures":{"signers":["0xa"],"signatures":["0xs"]}}}
		},
		"rate_data":{
			"BTC":{"1700000000":{"currency":"BTC","expiry":1700000000,"rate":"0.05","confidence":"0.95","timestamp":1700000000000,"deadline":1700000060000,"signatures":{"signers":["0xa"],"signatures":["0xs"]}}}
		},
		"vol_data":{
			"BTC":{"1700000000":{"currency":"BTC","expiry":1700000000,"vol_data":{"SVI_a":"0","SVI_b":"0.1","SVI_fwd":"50000","SVI_m":"0","SVI_refTau":"0.1","SVI_rho":"-0.5","SVI_sigma":"0.5"},"confidence":"0.95","timestamp":1700000000000,"deadline":1700000060000,"signatures":{"signers":["0xa"],"signatures":["0xs"]}}}
		}
	}`)
	var got types.SignedFeeds
	require.NoError(t, json.Unmarshal(raw, &got))

	require.Contains(t, got.SpotData, "BTC")
	assert.Equal(t, "50000", got.SpotData["BTC"].Price.String())
	assert.Equal(t, "S", got.SpotData["BTC"].FeedSourceType)

	require.Contains(t, got.PerpData, "BTC")
	require.Contains(t, got.PerpData["BTC"], "P")
	assert.Equal(t, "P", got.PerpData["BTC"]["P"].Type)
	assert.Equal(t, "100", got.PerpData["BTC"]["P"].SpotDiffValue.String())

	require.Contains(t, got.FwdData, "BTC")
	require.Contains(t, got.FwdData["BTC"], "1700000000")
	assert.Equal(t, "50", got.FwdData["BTC"]["1700000000"].FwdDiff.String())

	require.Contains(t, got.RateData, "BTC")
	assert.Equal(t, "0.05", got.RateData["BTC"]["1700000000"].Rate.String())

	require.Contains(t, got.VolData, "BTC")
	vol := got.VolData["BTC"]["1700000000"]
	assert.Equal(t, "0.5", vol.VolData.SVISigma.String())
	assert.Equal(t, "-0.5", vol.VolData.SVIRho.String())

	assert.Equal(t, []string{"0xa"}, got.SpotData["BTC"].Signatures.Signers)
	assert.Equal(t, []string{"0xs"}, got.SpotData["BTC"].Signatures.Signatures)
}

func TestVolSVIParam_RoundTrip(t *testing.T) {
	in := types.VolSVIParam{
		SVIA:      types.MustDecimal("0"),
		SVIB:      types.MustDecimal("0.1"),
		SVIFwd:    types.MustDecimal("50000"),
		SVIM:      types.MustDecimal("0"),
		SVIRefTau: types.MustDecimal("0.1"),
		SVIRho:    types.MustDecimal("-0.5"),
		SVISigma:  types.MustDecimal("0.5"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.VolSVIParam
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.SVIRho.String(), out.SVIRho.String())
	// JSON tag honours the OAS exact field name.
	assert.Contains(t, string(b), `"SVI_a"`)
	assert.Contains(t, string(b), `"SVI_refTau"`)
}

func TestOracleSignatureData_Empty(t *testing.T) {
	// Empty struct still marshals/unmarshals cleanly.
	in := types.OracleSignatureData{}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.OracleSignatureData
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Empty(t, out.Signers)
	assert.Empty(t, out.Signatures)
}
