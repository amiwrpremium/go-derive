package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestFundingPayment_Decode(t *testing.T) {
	raw := []byte(`{
		"instrument_name":"BTC-PERP",
		"subaccount_id":42,
		"timestamp":1700000000000,
		"funding":"0.000123",
		"pnl":"-0.00045"
	}`)
	var got types.FundingPayment
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "BTC-PERP", got.InstrumentName)
	assert.Equal(t, int64(42), got.SubaccountID)
	assert.Equal(t, int64(1700000000000), got.Timestamp.Millis())
	assert.Equal(t, "0.000123", got.Funding.String())
	assert.Equal(t, "-0.00045", got.PnL.String())
}

func TestInterestPayment_Decode(t *testing.T) {
	raw := []byte(`{"subaccount_id":1,"timestamp":1700000000000,"interest":"0.05"}`)
	var got types.InterestPayment
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "0.05", got.Interest.String())
}

func TestERC20Transfer_Decode(t *testing.T) {
	raw := []byte(`{
		"subaccount_id":1,
		"counterparty_subaccount_id":2,
		"asset":"USDC",
		"amount":"100",
		"is_outgoing":true,
		"timestamp":1700000000000,
		"tx_hash":"0x1111111111111111111111111111111111111111111111111111111111111111"
	}`)
	var got types.ERC20Transfer
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.True(t, got.IsOutgoing)
	assert.Equal(t, "USDC", got.Asset)
	assert.Equal(t, "100", got.Amount.String())
	assert.False(t, got.TxHash.IsZero())
}

func TestOptionSettlement_Decode(t *testing.T) {
	raw := []byte(`{
		"subaccount_id":1,
		"instrument_name":"BTC-20240101-50000-C",
		"expiry":1704067200,
		"amount":"1",
		"settlement_price":"50000",
		"option_settlement_pnl":"100",
		"option_settlement_pnl_excl_fees":"105"
	}`)
	var got types.OptionSettlement
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, int64(1704067200), got.Expiry)
	assert.Equal(t, "100", got.OptionSettlementPnL.String())
}

func TestSubaccountValueRecord_Decode(t *testing.T) {
	raw := []byte(`{
		"timestamp":1700000000000,
		"currency":"USDC",
		"margin_type":"PM",
		"subaccount_value":"10000",
		"initial_margin":"100",
		"maintenance_margin":"50",
		"delayed_maintenance_margin":"55"
	}`)
	var got types.SubaccountValueRecord
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "PM", got.MarginType)
	assert.Equal(t, "10000", got.SubaccountValue.String())
}

func TestExpiredAndCancelledExport_Decode(t *testing.T) {
	raw := []byte(`{"presigned_urls":["https://s3.example/a","https://s3.example/b"]}`)
	var got types.ExpiredAndCancelledExport
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, []string{"https://s3.example/a", "https://s3.example/b"}, got.PresignedURLs)
}
