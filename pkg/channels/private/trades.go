// Package private declares Derive's authenticated WebSocket subscription
// channels: a subaccount's order, position, balance, trade, RFQ and quote
// streams.
//
// Each descriptor needs a SubaccountID and the WebSocket session must be
// logged in (call [github.com/amiwrpremium/go-derive/pkg/ws.Client.Login]
// before subscribing). Pair them with the matching T when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe].
package private

import (
	"encoding/json"
	"fmt"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// Trades subscribes to fill events for one subaccount.
//
// The dotted server-side channel name (per
// https://docs.derive.xyz/reference/subaccount_id-trades) is:
//
//	{subaccount_id}.trades
//
// Pair this descriptor with T = [[]types.Trade] when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe]. Each notification
// is a batch of fills since the last update.
type Trades struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
}

// Name returns the dotted server-side channel string.
func (t Trades) Name() string { return fmt.Sprintf("%d.trades", t.SubaccountID) }

// Decode parses an inbound notification payload into a [[]types.Trade].
func (Trades) Decode(raw json.RawMessage) (any, error) {
	var trades []types.Trade
	if err := json.Unmarshal(raw, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}

// TradesByTxStatus subscribes to fill events for one subaccount, filtered
// server-side by on-chain transaction status. Useful for makers who only
// want to see settled fills.
//
// The dotted server-side channel name (per
// https://docs.derive.xyz/reference/subaccount_id-trades-tx_status) is:
//
//	{subaccount_id}.trades.{tx_status}
//
// Per the docs, only `settled` and `reverted` are documented filter
// values today — other [enums.TxStatus] values may be rejected by the
// engine. Same per-batch payload as [Trades]; pair with T = [[]types.Trade].
type TradesByTxStatus struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
	// TxStatus is the on-chain transaction status to filter by.
	// Documented values: [enums.TxStatusSettled], [enums.TxStatusReverted].
	TxStatus enums.TxStatus
}

// Name returns the dotted server-side channel string.
func (t TradesByTxStatus) Name() string {
	return fmt.Sprintf("%d.trades.%s", t.SubaccountID, t.TxStatus)
}

// Decode parses an inbound notification payload into a [[]types.Trade].
// Same payload shape as the unfiltered [Trades] channel.
func (TradesByTxStatus) Decode(raw json.RawMessage) (any, error) {
	var trades []types.Trade
	if err := json.Unmarshal(raw, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}
