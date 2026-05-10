package public_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/channels/public"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestAuctionsWatch_Name(t *testing.T) {
	assert.Equal(t, "auctions.watch", public.AuctionsWatch{}.Name())
}

func TestAuctionsWatch_Decode_Ongoing(t *testing.T) {
	raw := json.RawMessage(`{"subaccount_id":7,"state":"ongoing",` +
		`"timestamp":1700000000000,"details":{"currency":"USDC",` +
		`"estimated_bid_price":"-12.5","estimated_discount_pnl":"1.2",` +
		`"estimated_mtm":"100","estimated_percent_bid":"0.25",` +
		`"last_seen_trade_id":42,"margin_type":"PM",` +
		`"min_cash_transfer":"50","min_price_limit":"5",` +
		`"subaccount_balances":{"USDC":"100"}}}`)
	got, err := public.AuctionsWatch{}.Decode(raw)
	require.NoError(t, err)
	ev, ok := got.(types.AuctionWatchEvent)
	require.True(t, ok, "Decode must return types.AuctionWatchEvent")
	assert.Equal(t, int64(7), ev.SubaccountID)
	assert.Equal(t, enums.AuctionStateOngoing, ev.State)
	require.NotNil(t, ev.Details)
	assert.Equal(t, "USDC", ev.Details.Currency)
}

func TestAuctionsWatch_Decode_Ended(t *testing.T) {
	raw := json.RawMessage(`{"subaccount_id":9,"state":"ended",` +
		`"timestamp":1700000060000,"details":null}`)
	got, err := public.AuctionsWatch{}.Decode(raw)
	require.NoError(t, err)
	ev := got.(types.AuctionWatchEvent)
	assert.Equal(t, enums.AuctionStateEnded, ev.State)
	assert.Nil(t, ev.Details)
}
