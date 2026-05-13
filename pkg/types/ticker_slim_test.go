package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestTickerSlim_DecodeFull(t *testing.T) {
	// Sample mirroring the real Derive `ticker_slim.<inst>.<interval>`
	// notification: outer envelope + abbreviated inner fields.
	payload := `{
		"timestamp": 1700000001000,
		"instrument_ticker": {
			"t": 1700000001000,
			"A": "0.5",  "a": "78758.5",
			"B": "0.4",  "b": "78752.1",
			"I": "78760", "M": "78755",
			"minp": "78000", "maxp": "79500",
			"f": "0.0001",
			"stats": {"v": "100"},
			"option_pricing": null
		}
	}`
	var ts types.TickerSlim
	require.NoError(t, json.Unmarshal([]byte(payload), &ts))
	assert.Equal(t, int64(1700000001000), ts.Timestamp.Millis())
	assert.Equal(t, int64(1700000001000), ts.Ticker.Timestamp.Millis())
	assert.Equal(t, "0.5", ts.Ticker.BestAskAmount.String())
	assert.Equal(t, "78758.5", ts.Ticker.BestAskPrice.String())
	assert.Equal(t, "0.4", ts.Ticker.BestBidAmount.String())
	assert.Equal(t, "78752.1", ts.Ticker.BestBidPrice.String())
	assert.Equal(t, "78760", ts.Ticker.IndexPrice.String())
	assert.Equal(t, "78755", ts.Ticker.MarkPrice.String())
	assert.Equal(t, "78000", ts.Ticker.MinPrice.String())
	assert.Equal(t, "79500", ts.Ticker.MaxPrice.String())
	assert.Equal(t, "0.0001", ts.Ticker.FundingRate.String())
	assert.JSONEq(t, `{"v":"100"}`, string(ts.Ticker.Stats))
	// JSON `null` decodes into RawMessage as the literal bytes `null`,
	// not Go nil — keep the field type future-proof for non-null option blocks.
	assert.Equal(t, "null", string(ts.Ticker.OptionPricing))
}

func TestTickerSlim_DecodeMinimal(t *testing.T) {
	// Only required fields populated; optional Decimals + RawMessages absent.
	payload := `{
		"timestamp": 1,
		"instrument_ticker": {
			"t": 1,
			"A": "0", "a": "0",
			"B": "0", "b": "0"
		}
	}`
	var ts types.TickerSlim
	require.NoError(t, json.Unmarshal([]byte(payload), &ts))
	assert.Equal(t, "0", ts.Ticker.BestAskAmount.String())
	assert.Equal(t, json.RawMessage(nil), ts.Ticker.Stats)
	assert.Equal(t, json.RawMessage(nil), ts.Ticker.OptionPricing)
}

func TestTickerSlim_RoundTrip_PreservesTopFields(t *testing.T) {
	var in types.TickerSlim
	in.Ticker.MarkPrice = types.MustDecimal("100.5")
	in.Ticker.IndexPrice = types.MustDecimal("100.5")
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.TickerSlim
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.Ticker.MarkPrice.String(), out.Ticker.MarkPrice.String())
}
