package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestInstrumentTickerFeed_Decode(t *testing.T) {
	raw := []byte(`{
		"timestamp": 1700000000123,
		"instrument_ticker": {
			"instrument_name": "BTC-PERP",
			"instrument_type": "perp",
			"is_active": true,
			"base_currency": "BTC",
			"quote_currency": "USD",
			"base_asset_address": "0x1111111111111111111111111111111111111111",
			"base_asset_sub_id": "0",
			"amount_step": "0.001",
			"minimum_amount": "0.001",
			"maximum_amount": "100",
			"tick_size": "0.5",
			"base_fee": "0.1",
			"maker_fee_rate": "0.0003",
			"taker_fee_rate": "0.0005",
			"mark_price_fee_rate_cap": null,
			"scheduled_activation": 1700000000,
			"scheduled_deactivation": 9223372036854775807,
			"best_bid_price": "65000",
			"best_bid_amount": "1.5",
			"best_ask_price": "65010",
			"best_ask_amount": "2.0",
			"mark_price": "65005",
			"index_price": "65003",
			"min_price": "64000",
			"max_price": "66000",
			"timestamp": 1700000000122,
			"option_details": null,
			"perp_details": {"index": "BTC", "max_leverage": "20"},
			"option_pricing": null,
			"stats": {"contract_volume":"100","high":"65500","low":"64500","num_trades":"42","open_interest":"10","percent_change":"0.01","usd_change":"650"}
		}
	}`)
	var f types.InstrumentTickerFeed
	require.NoError(t, json.Unmarshal(raw, &f))
	assert.Equal(t, int64(1700000000123), f.Timestamp.Millis())
	assert.Equal(t, "BTC-PERP", f.Ticker.InstrumentName)
	assert.Equal(t, enums.InstrumentTypePerp, f.Ticker.InstrumentType)
	assert.True(t, f.Ticker.IsActive)
	assert.Equal(t, "65005", f.Ticker.MarkPrice.String())
	assert.Equal(t, "0.5", f.Ticker.TickSize.String())
	assert.Equal(t, "0", f.Ticker.MarkPriceFeeRateCap.String(), "null decodes to zero")
	assert.NotEmpty(t, f.Ticker.PerpDetails, "perp_details preserved as raw JSON")
	assert.NotEmpty(t, f.Ticker.Stats, "stats preserved as raw JSON")
}

func TestInstrumentTickerFeed_Decode_WithMatchingFieldsAndDepth(t *testing.T) {
	// Ticker channel publishes the matching-engine parameters
	// (pro_rata_*, fifo_min_allocation) and 5%-depth liquidity
	// markers (five_percent_*_depth). All five are nullable, so
	// the channel may omit them entirely on some instruments.
	raw := []byte(`{
		"timestamp": 1700000000123,
		"instrument_ticker": {
			"instrument_name": "BTC-PERP",
			"instrument_type": "perp",
			"is_active": true,
			"base_currency": "BTC",
			"quote_currency": "USD",
			"base_asset_address": "0x1111111111111111111111111111111111111111",
			"base_asset_sub_id": "0",
			"amount_step": "0.001",
			"minimum_amount": "0.001",
			"maximum_amount": "100",
			"tick_size": "0.5",
			"base_fee": "0.1",
			"maker_fee_rate": "0.0003",
			"taker_fee_rate": "0.0005",
			"pro_rata_fraction": "0.5",
			"pro_rata_amount_step": "0.01",
			"fifo_min_allocation": "0.1",
			"scheduled_activation": 1700000000,
			"scheduled_deactivation": 9223372036854775807,
			"best_bid_price": "65000",
			"best_bid_amount": "1.5",
			"best_ask_price": "65010",
			"best_ask_amount": "2.0",
			"five_percent_bid_depth": "12.5",
			"five_percent_ask_depth": "13.25",
			"mark_price": "65005",
			"index_price": "65003",
			"min_price": "64000",
			"max_price": "66000",
			"timestamp": 1700000000122,
			"stats": {}
		}
	}`)
	var f types.InstrumentTickerFeed
	require.NoError(t, json.Unmarshal(raw, &f))
	assert.Equal(t, "0.5", f.Ticker.ProRataFraction.String())
	assert.Equal(t, "0.01", f.Ticker.ProRataAmountStep.String())
	assert.Equal(t, "0.1", f.Ticker.FIFOMinAllocation.String())
	assert.Equal(t, "12.5", f.Ticker.FivePercentBidDepth.String())
	assert.Equal(t, "13.25", f.Ticker.FivePercentAskDepth.String())
}
