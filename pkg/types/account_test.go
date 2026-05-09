package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestAccountResult_Decode(t *testing.T) {
	raw := []byte(`{
		"subaccount_ids":[1,2],
		"wallet":"0xabcd",
		"cancel_on_disconnect":true,
		"creation_timestamp_sec":1700000000,
		"is_rfq_maker":false,
		"referral_code":"REF-123",
		"websocket_matching_tps":50,
		"websocket_non_matching_tps":20,
		"websocket_option_tps":10,
		"websocket_perp_tps":15,
		"per_endpoint_tps":{"public/get_time":100},
		"fee_info":{
			"base_fee_discount":"0.1",
			"option_maker_fee":"0.0003",
			"option_taker_fee":"0.0005",
			"perp_maker_fee":"0.0001",
			"perp_taker_fee":"0.0003",
			"rfq_maker_discount":"0.5",
			"rfq_taker_discount":"0.5",
			"spot_maker_fee":"0",
			"spot_taker_fee":"0.0002"
		}
	}`)

	var got types.AccountResult
	require.NoError(t, json.Unmarshal(raw, &got))

	assert.Equal(t, []int64{1, 2}, got.SubaccountIDs)
	assert.Equal(t, "0xabcd", got.Wallet)
	assert.True(t, got.CancelOnDisconnect)
	assert.Equal(t, int64(1700000000), got.CreationTimestampSec)
	assert.False(t, got.IsRFQMaker)
	assert.Equal(t, "REF-123", got.ReferralCode)
	assert.Equal(t, int64(50), got.WebSocketMatchingTPS)
	assert.Equal(t, int64(20), got.WebSocketNonMatchingTPS)
	assert.Equal(t, int64(10), got.WebSocketOptionTPS)
	assert.Equal(t, int64(15), got.WebSocketPerpTPS)
	assert.JSONEq(t, `{"public/get_time":100}`, string(got.PerEndpointTPS))
	assert.Equal(t, "0.0003", got.FeeInfo.OptionMakerFee.String())
	assert.Equal(t, "0.5", got.FeeInfo.RFQMakerDiscount.String())
}

func TestMarginResult_Decode(t *testing.T) {
	raw := []byte(`{
		"subaccount_id":42,
		"is_valid_trade":true,
		"pre_initial_margin":"100.5",
		"post_initial_margin":"110.25",
		"pre_maintenance_margin":"50",
		"post_maintenance_margin":"55"
	}`)

	var got types.MarginResult
	require.NoError(t, json.Unmarshal(raw, &got))

	assert.Equal(t, int64(42), got.SubaccountID)
	assert.True(t, got.IsValidTrade)
	assert.Equal(t, "100.5", got.PreInitialMargin.String())
	assert.Equal(t, "110.25", got.PostInitialMargin.String())
	assert.Equal(t, "50", got.PreMaintenanceMargin.String())
	assert.Equal(t, "55", got.PostMaintenanceMargin.String())
}

func TestFeeInfo_RoundTrip(t *testing.T) {
	in := types.FeeInfo{
		BaseFeeDiscount:  types.MustDecimal("0.1"),
		OptionMakerFee:   types.MustDecimal("0.0003"),
		OptionTakerFee:   types.MustDecimal("0.0005"),
		PerpMakerFee:     types.MustDecimal("0.0001"),
		PerpTakerFee:     types.MustDecimal("0.0003"),
		RFQMakerDiscount: types.MustDecimal("0.5"),
		RFQTakerDiscount: types.MustDecimal("0.5"),
		SpotMakerFee:     types.MustDecimal("0"),
		SpotTakerFee:     types.MustDecimal("0.0002"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.FeeInfo
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.OptionMakerFee.String(), out.OptionMakerFee.String())
	assert.Equal(t, in.SpotMakerFee.String(), out.SpotMakerFee.String())
}
