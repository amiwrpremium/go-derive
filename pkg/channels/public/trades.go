// Package public declares Derive's public (no-auth) WebSocket subscription
// channels: order books, public trade prints, tickers, and instrument
// add/remove events.
//
// Every descriptor in this package satisfies
// [github.com/amiwrpremium/go-derive/pkg/channels.Channel]; pass them to
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe] together with a
// matching T.
package public

import (
	"encoding/json"
	"fmt"

	"github.com/amiwrpremium/go-derive"
)

// Trades subscribes to public trade prints on one instrument.
//
// The dotted server-side channel name is:
//
//	trades.{instrument}
//
// Pair this descriptor with T = []derive.Trade when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe]. Each notification
// carries a batch of [derive.Trade] events that printed in the same window.
type Trades struct {
	// Instrument is the market name (e.g. "BTC-PERP").
	Instrument string
}

// Name returns the dotted server-side channel string.
func (t Trades) Name() string { return fmt.Sprintf("trades.%s", t.Instrument) }

// Decode parses an inbound notification payload into a [[]derive.Trade].
func (Trades) Decode(raw json.RawMessage) (any, error) {
	var trades []derive.Trade
	if err := json.Unmarshal(raw, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}
