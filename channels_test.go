package derive_test

import (
	"encoding/json"
	"github.com/amiwrpremium/go-derive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// stubChannel is a minimal in-test implementation used only to confirm the
// Channel interface is satisfiable and that pkg/ws.Subscribe works against
// arbitrary descriptors. Real descriptors live in pkg/channels/{public,private}.
type stubChannel struct{ name string }

func (s stubChannel) Name() string { return s.name }
func (stubChannel) Decode(raw json.RawMessage) (any, error) {
	var out string
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func TestChannel_InterfaceConformance(t *testing.T) {
	var c derive.Channel = stubChannel{name: "trades.BTC-PERP"}
	assert.Equal(t, "trades.BTC-PERP", c.Name())

	v, err := c.Decode(json.RawMessage(`"hello"`))
	require.NoError(t, err)
	assert.Equal(t, "hello", v)
}

func TestChannel_DecodeError(t *testing.T) {
	c := stubChannel{name: "x"}
	_, err := c.Decode(json.RawMessage(`{`))
	assert.Error(t, err)
}
func TestOrderBook_Name_DefaultsApplied(t *testing.T) {
	got := derive.PublicOrderBook{Instrument: "BTC-PERP"}.Name()
	assert.Equal(t, "orderbook.BTC-PERP.1.10", got)
}

func TestOrderBook_Name_ExplicitGroupAndDepth(t *testing.T) {
	got := derive.PublicOrderBook{Instrument: "ETH-PERP", Group: "10", Depth: 25}.Name()
	assert.Equal(t, "orderbook.ETH-PERP.10.25", got)
}

func TestOrderBook_Name_ExplicitGroupDefaultDepth(t *testing.T) {
	got := derive.PublicOrderBook{Instrument: "X", Group: "5", Depth: 0}.Name()
	assert.Equal(t, "orderbook.X.5.10", got)
	assert.Equal(t, "orderbook.X.5.10", derive.PublicOrderBook{Instrument: "X", Group: "5"}.Name())
}

func TestOrderBook_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`{
		"instrument_name":"BTC-PERP",
		"bids":[["100","1"]],
		"asks":[["101","2"]],
		"timestamp":1700000000000
	}`)
	v, err := derive.PublicOrderBook{}.Decode(raw)
	require.NoError(t, err)
	ob, ok := v.(derive.OrderBook)
	require.True(t, ok)
	assert.Equal(t, "BTC-PERP", ob.InstrumentName)
}

func TestOrderBook_Decode_Malformed(t *testing.T) {
	_, err := derive.PublicOrderBook{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
func TestSpotFeed_Name(t *testing.T) {
	assert.Equal(t, "spot_feed.BTC", derive.PublicSpotFeed{Currency: "BTC"}.Name())
	assert.Equal(t, "spot_feed.ETH", derive.PublicSpotFeed{Currency: "ETH"}.Name())
}

func TestPublicSpotFeed_Decode(t *testing.T) {
	raw := json.RawMessage(`{"timestamp":1,"feeds":{"BTC":{"price":"100","confidence":"1","price_prev_daily":"99","confidence_prev_daily":"1","timestamp_prev_daily":0}}}`)
	v, err := derive.PublicSpotFeed{}.Decode(raw)
	require.NoError(t, err)
	sf, ok := v.(derive.SpotFeed)
	require.True(t, ok)
	assert.Equal(t, "100", sf.Feeds["BTC"].Price.String())
}

func TestPublicSpotFeed_Decode_Malformed(t *testing.T) {
	_, err := derive.PublicSpotFeed{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
func TestTickerSlim_Name_DefaultInterval(t *testing.T) {
	got := derive.PublicTickerSlim{Instrument: "BTC-PERP"}.Name()
	assert.Equal(t, "ticker_slim.BTC-PERP.1000", got)
}

func TestTickerSlim_Name_ExplicitInterval(t *testing.T) {
	got := derive.PublicTickerSlim{Instrument: "BTC-PERP", Interval: "100"}.Name()
	assert.Equal(t, "ticker_slim.BTC-PERP.100", got)
}

func TestTickerSlim_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`{
		"timestamp": 1700000000000,
		"instrument_ticker": {
			"t": 1700000000000,
			"A": "0.5", "a": "78758.5",
			"B": "0.4", "b": "78752.1",
			"M": "78755", "I": "78760",
			"f": "0.0001"
		}
	}`)
	v, err := derive.PublicTickerSlim{}.Decode(raw)
	require.NoError(t, err)
	tk, ok := v.(derive.TickerSlim)
	require.True(t, ok)
	assert.Equal(t, "78752.1", tk.Ticker.BestBidPrice.String())
	assert.Equal(t, "78758.5", tk.Ticker.BestAskPrice.String())
	assert.Equal(t, "78755", tk.Ticker.MarkPrice.String())
}

func TestTickerSlim_Decode_Malformed(t *testing.T) {
	_, err := derive.PublicTickerSlim{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
func TestTradesByType_Name(t *testing.T) {
	assert.Equal(t, "trades.perp.BTC",
		derive.PublicTradesByType{InstrumentType: derive.InstrumentTypePerp, Currency: "BTC"}.Name())
	assert.Equal(t, "trades.option.ETH",
		derive.PublicTradesByType{InstrumentType: derive.InstrumentTypeOption, Currency: "ETH"}.Name())
}

func TestTradesByType_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"trade_id":"t1","instrument_name":"BTC-PERP","direction":"buy","trade_price":"100","trade_amount":"1","timestamp":1700000000000}]`)
	v, err := derive.PublicTradesByType{}.Decode(raw)
	require.NoError(t, err)
	trades, ok := v.([]derive.Trade)
	require.True(t, ok)
	require.Len(t, trades, 1)
	assert.Equal(t, "t1", trades[0].TradeID)
}

func TestTradesByType_Decode_Malformed(t *testing.T) {
	_, err := derive.PublicTradesByType{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
func TestTrades_Name(t *testing.T) {
	assert.Equal(t, "trades.BTC-PERP", derive.PublicTrades{Instrument: "BTC-PERP"}.Name())
	assert.Equal(t, "trades.", derive.PublicTrades{}.Name())
}

func TestTrades_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"trade_id":"T1","instrument_name":"BTC-PERP","direction":"buy","trade_price":"65000","trade_amount":"0.1","mark_price":"65000","timestamp":1700000000000}]`)
	v, err := derive.PublicTrades{}.Decode(raw)
	require.NoError(t, err)
	trades, ok := v.([]derive.Trade)
	require.True(t, ok)
	require.Len(t, trades, 1)
	assert.Equal(t, "T1", trades[0].TradeID)
}

func TestTrades_Decode_EmptyArray(t *testing.T) {
	v, err := derive.PublicTrades{}.Decode(json.RawMessage(`[]`))
	require.NoError(t, err)
	trades := v.([]derive.Trade)
	assert.Empty(t, trades)
}

func TestTrades_Decode_Malformed(t *testing.T) {
	_, err := derive.PublicTrades{}.Decode([]byte(`not-json`))
	assert.Error(t, err)
}
func TestBalances_Name(t *testing.T) {
	assert.Equal(t, "subaccount.5.balances", derive.PrivateBalances{SubaccountID: 5}.Name())
}

func TestBalances_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`{"subaccount_id":5,"subaccount_value":"100","initial_margin":"50","maintenance_margin":"30","collaterals":[]}`)
	v, err := derive.PrivateBalances{}.Decode(raw)
	require.NoError(t, err)
	bal, ok := v.(derive.Balance)
	require.True(t, ok)
	assert.Equal(t, int64(5), bal.SubaccountID)
}

func TestBalances_Decode_Malformed(t *testing.T) {
	_, err := derive.PrivateBalances{}.Decode([]byte(`[`))
	assert.Error(t, err)
}
func TestOrders_Name(t *testing.T) {
	assert.Equal(t, "subaccount.123.orders", derive.PrivateOrders{SubaccountID: 123}.Name())
}

func TestOrders_Name_Zero(t *testing.T) {
	assert.Equal(t, "subaccount.0.orders", derive.PrivateOrders{}.Name())
}

func TestOrders_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"order_id":"O1","subaccount_id":1,"instrument_name":"BTC-PERP","direction":"buy","order_type":"limit","time_in_force":"gtc","order_status":"open","amount":"0.1","filled_amount":"0","limit_price":"65000","max_fee":"10","nonce":1,"signer":"0x0000000000000000000000000000000000000000","creation_timestamp":1700000000000,"last_update_timestamp":1700000000000}]`)
	v, err := derive.PrivateOrders{}.Decode(raw)
	require.NoError(t, err)
	orders, ok := v.([]derive.Order)
	require.True(t, ok)
	require.Len(t, orders, 1)
	assert.Equal(t, "O1", orders[0].OrderID)
}

func TestOrders_Decode_Malformed(t *testing.T) {
	_, err := derive.PrivateOrders{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
func TestRFQs_Name(t *testing.T) {
	assert.Equal(t,
		"wallet.0xC0FFee0000000000000000000000000000000000.rfqs",
		derive.PrivateRFQs{Wallet: "0xC0FFee0000000000000000000000000000000000"}.Name())
}

func TestRFQs_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"rfq_id":"R","subaccount_id":9,"status":"open","legs":[],"creation_timestamp":1,"last_update_timestamp":1}]`)
	v, err := derive.PrivateRFQs{}.Decode(raw)
	require.NoError(t, err)
	rfqs, ok := v.([]derive.RFQ)
	require.True(t, ok)
	require.Len(t, rfqs, 1)
	assert.Equal(t, "R", rfqs[0].RFQID)
}

func TestRFQs_Decode_Malformed(t *testing.T) {
	_, err := derive.PrivateRFQs{}.Decode([]byte(`{`))
	assert.Error(t, err)
}

func TestQuotes_Name(t *testing.T) {
	assert.Equal(t, "subaccount.9.quotes", derive.PrivateQuotes{SubaccountID: 9}.Name())
}

func TestQuotes_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"quote_id":"Q","rfq_id":"R","subaccount_id":9,"direction":"buy","legs":[],"price":"1","status":"open","creation_timestamp":1}]`)
	v, err := derive.PrivateQuotes{}.Decode(raw)
	require.NoError(t, err)
	quotes, ok := v.([]derive.Quote)
	require.True(t, ok)
	require.Len(t, quotes, 1)
}

func TestQuotes_Decode_Malformed(t *testing.T) {
	_, err := derive.PrivateQuotes{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
func TestPrivateTrades_Name(t *testing.T) {
	assert.Equal(t, "subaccount.7.trades", derive.PrivateTrades{SubaccountID: 7}.Name())
}

func TestPrivateTrades_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"trade_id":"T","instrument_name":"BTC-PERP","direction":"buy","trade_price":"1","trade_amount":"1","mark_price":"1","timestamp":1700000000000}]`)
	v, err := derive.PrivateTrades{}.Decode(raw)
	require.NoError(t, err)
	trades, ok := v.([]derive.Trade)
	require.True(t, ok)
	require.Len(t, trades, 1)
}

func TestPrivateTrades_Decode_Malformed(t *testing.T) {
	_, err := derive.PrivateTrades{}.Decode([]byte(`{`))
	assert.Error(t, err)
}
