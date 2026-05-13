package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestTicker_DecodeFull(t *testing.T) {
	// Mirrors the real Derive response shape: per-margin-type open_interest
	// breakdown, 5%-depth fields, and price bands.
	payload := `{
		"instrument_name": "BTC-PERP",
		"instrument_type": "perp",
		"is_active": true,
		"best_bid_price": "100",
		"best_bid_amount": "1",
		"best_ask_price": "101",
		"best_ask_amount": "2",
		"five_percent_bid_depth": "10",
		"five_percent_ask_depth": "12",
		"mark_price": "100.5",
		"index_price": "100.5",
		"min_price": "90",
		"max_price": "110",
		"open_interest": {"PM": [{"current_open_interest": "1", "interest_cap": "100", "manager_currency": "BTC"}]},
		"timestamp": 1700000000000
	}`
	var tk types.Ticker
	require.NoError(t, json.Unmarshal([]byte(payload), &tk))
	assert.Equal(t, "BTC-PERP", tk.InstrumentName)
	assert.Equal(t, "perp", tk.InstrumentType)
	assert.True(t, tk.IsActive)
	assert.Equal(t, "100", tk.BestBidPrice.String())
	assert.Equal(t, "10", tk.FivePercentBidDepth.String())
	assert.Equal(t, "90", tk.MinPrice.String())
	assert.Contains(t, string(tk.OpenInterest), "current_open_interest")
}

func TestTicker_DecodeMinimal(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"best_bid_price": "0",
		"best_bid_amount": "0",
		"best_ask_price": "0",
		"best_ask_amount": "0",
		"mark_price": "0",
		"index_price": "0",
		"timestamp": 0
	}`
	var tk types.Ticker
	require.NoError(t, json.Unmarshal([]byte(payload), &tk))
	assert.Equal(t, "BTC-PERP", tk.InstrumentName)
}

func TestTicker_DecodeInstrumentMetadata(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"instrument_type": "perp",
		"is_active": true,
		"base_currency": "BTC",
		"quote_currency": "USDC",
		"base_asset_address": "0x1111111111111111111111111111111111111111",
		"base_asset_sub_id": "0",
		"tick_size": "0.5",
		"amount_step": "0.001",
		"minimum_amount": "0.001",
		"maximum_amount": "1000",
		"maker_fee_rate": "0.0001",
		"taker_fee_rate": "0.0005",
		"base_fee": "0.5",
		"mark_price_fee_rate_cap": "0.1",
		"pro_rata_fraction": "0",
		"pro_rata_amount_step": "0",
		"fifo_min_allocation": "0",
		"scheduled_activation": 1700000000,
		"scheduled_deactivation": 9223372036854775807,
		"best_bid_price": "100",
		"best_bid_amount": "1",
		"best_ask_price": "101",
		"best_ask_amount": "2",
		"mark_price": "100.5",
		"index_price": "100.5",
		"stats": {"daily_volume": "1000"},
		"option_pricing": null,
		"perp_details": {"index": "BTC", "max_leverage": "50"},
		"timestamp": 1700000000000
	}`
	var tk types.Ticker
	require.NoError(t, json.Unmarshal([]byte(payload), &tk))
	assert.Equal(t, "BTC", tk.BaseCurrency)
	assert.Equal(t, "USDC", tk.QuoteCurrency)
	assert.Equal(t, "0.5", tk.TickSize.String())
	assert.Equal(t, "0.001", tk.AmountStep.String())
	assert.Equal(t, "0.0001", tk.MakerFeeRate.String())
	assert.Equal(t, "0.0005", tk.TakerFeeRate.String())
	assert.Equal(t, int64(1700000000), tk.ScheduledActivation)
	require.NotNil(t, tk.Perp)
	assert.Equal(t, "BTC", tk.Perp.IndexName)
	assert.Contains(t, string(tk.Stats), "daily_volume")
}

func TestTicker_RoundTrip(t *testing.T) {
	in := types.Ticker{
		InstrumentName: "BTC-PERP",
		BestBidPrice:   types.MustDecimal("100"),
		BestAskPrice:   types.MustDecimal("101"),
		MarkPrice:      types.MustDecimal("100.5"),
		IndexPrice:     types.MustDecimal("100.5"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.Ticker
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.InstrumentName, out.InstrumentName)
	assert.Equal(t, in.BestBidPrice.String(), out.BestBidPrice.String())
}
