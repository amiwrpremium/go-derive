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
