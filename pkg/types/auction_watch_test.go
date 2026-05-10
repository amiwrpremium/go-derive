package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestAuctionWatchEvent_Decode_Ongoing(t *testing.T) {
	raw := []byte(`{"subaccount_id":7,"state":"ongoing","timestamp":1700000000000,` +
		`"details":{"currency":"USDC","estimated_bid_price":"-12.5",` +
		`"estimated_discount_pnl":"1.2","estimated_mtm":"100",` +
		`"estimated_percent_bid":"0.25","last_seen_trade_id":42,` +
		`"margin_type":"PM","min_cash_transfer":"50","min_price_limit":"5",` +
		`"subaccount_balances":{"USDC":"100"}}}`)
	var ev types.AuctionWatchEvent
	require.NoError(t, json.Unmarshal(raw, &ev))
	assert.Equal(t, int64(7), ev.SubaccountID)
	assert.Equal(t, enums.AuctionStateOngoing, ev.State)
	assert.Equal(t, int64(1700000000000), ev.Timestamp.Millis())
	require.NotNil(t, ev.Details)
	assert.Equal(t, "USDC", ev.Details.Currency)
	assert.Equal(t, "-12.5", ev.Details.EstimatedBidPrice.String())
	assert.Equal(t, enums.MarginTypePM, ev.Details.MarginType)
	assert.Equal(t, int64(42), ev.Details.LastSeenTradeID)
	assert.JSONEq(t, `{"USDC":"100"}`, string(ev.Details.SubaccountBalances))
}

func TestAuctionWatchEvent_Decode_Ended(t *testing.T) {
	raw := []byte(`{"subaccount_id":9,"state":"ended","timestamp":1700000060000,"details":null}`)
	var ev types.AuctionWatchEvent
	require.NoError(t, json.Unmarshal(raw, &ev))
	assert.Equal(t, enums.AuctionStateEnded, ev.State)
	assert.Nil(t, ev.Details)
}
