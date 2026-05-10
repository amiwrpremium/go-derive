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
