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

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// Trades subscribes to fill events for one subaccount.
//
// The dotted server-side channel name is:
//
//	subaccount.{id}.trades
//
// Pair this descriptor with T = [[]types.Trade] when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe]. Each notification
// is a batch of fills since the last update.
type Trades struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
}

// Name returns the dotted server-side channel string.
func (t Trades) Name() string { return fmt.Sprintf("subaccount.%d.trades", t.SubaccountID) }

// Decode parses an inbound notification payload into a [[]types.Trade].
func (Trades) Decode(raw json.RawMessage) (any, error) {
	var trades []types.Trade
	if err := json.Unmarshal(raw, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}
