package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestInstrument_DecodePerp(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"base_currency": "BTC",
		"quote_currency": "USDC",
		"instrument_type": "perp",
		"is_active": true,
		"tick_size": "0.5",
		"minimum_amount": "0.001",
		"maximum_amount": "1000",
		"amount_step": "0.001",
		"mark_price": "65000.5",
		"index_price": "65000",
		"perp_details": {"index": "BTC", "max_leverage": "50"}
	}`
	var inst types.Instrument
	require.NoError(t, json.Unmarshal([]byte(payload), &inst))
	assert.Equal(t, derive.InstrumentTypePerp, inst.Type)
	require.NotNil(t, inst.Perp)
	assert.Equal(t, "BTC", inst.Perp.IndexName)
	assert.Equal(t, "50", inst.Perp.MaxLeverage.String())
	assert.Nil(t, inst.Option)
	assert.Nil(t, inst.ERC20)
}

func TestInstrument_DecodeOption(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-25DEC25-65000-C",
		"instrument_type": "option",
		"is_active": true,
		"tick_size": "0.5",
		"minimum_amount": "0.01",
		"maximum_amount": "100",
		"amount_step": "0.01",
		"option_details": {"option_type": "C", "strike": "65000", "expiry": 1735689600000, "index": "BTC"}
	}`
	var inst types.Instrument
	require.NoError(t, json.Unmarshal([]byte(payload), &inst))
	require.NotNil(t, inst.Option)
	assert.Equal(t, derive.OptionTypeCall, inst.Option.OptionType)
	assert.Equal(t, "65000", inst.Option.Strike.String())
}

func TestInstrument_DecodeERC20(t *testing.T) {
	payload := `{
		"instrument_name": "USDC",
		"instrument_type": "erc20",
		"is_active": true,
		"tick_size": "0.01",
		"minimum_amount": "1",
		"maximum_amount": "1000000",
		"amount_step": "1",
		"erc20_details": {
			"underlying_erc20_address": "0x1111111111111111111111111111111111111111",
			"borrow_index": "1.0",
			"supply_index": "1.0"
		}
	}`
	var inst types.Instrument
	require.NoError(t, json.Unmarshal([]byte(payload), &inst))
	require.NotNil(t, inst.ERC20)
	// shopspring normalises "1.0" to "1".
	assert.Equal(t, "1", inst.ERC20.BorrowIndex.String())
}

func TestInstrument_RoundTrip(t *testing.T) {
	in := types.Instrument{
		Name:          "BTC-PERP",
		BaseCurrency:  "BTC",
		QuoteCurrency: "USDC",
		Type:          derive.InstrumentTypePerp,
		IsActive:      true,
		TickSize:      types.MustDecimal("0.5"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.Instrument
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.Name, out.Name)
	assert.Equal(t, in.Type, out.Type)
}
