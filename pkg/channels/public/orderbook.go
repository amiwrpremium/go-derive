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

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// OrderBook subscribes to incremental order-book updates for one instrument.
//
// The dotted server-side channel name is:
//
//	orderbook.{instrument}.{group}.{depth}
//
// where Group is the price grouping (e.g. "10" for nearest-10-ticks; "1"
// for no grouping) and Depth is the number of levels returned per side.
// Empty Group defaults to "1"; zero Depth defaults to 10.
//
// Pair this descriptor with T = [types.OrderBook] when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe].
type OrderBook struct {
	// Instrument is the market name (e.g. "BTC-PERP").
	Instrument string
	// Group is the price-grouping string. Empty means "1" (no grouping).
	Group string
	// Depth is the number of book levels per side. Zero means 10.
	Depth int
}

// Name returns the dotted server-side channel string. See [OrderBook] for
// the format.
func (o OrderBook) Name() string {
	g := o.Group
	if g == "" {
		g = "1"
	}
	d := o.Depth
	if d == 0 {
		d = 10
	}
	return fmt.Sprintf("orderbook.%s.%s.%d", o.Instrument, g, d)
}

// Decode parses one inbound notification payload into a [types.OrderBook].
func (OrderBook) Decode(raw json.RawMessage) (any, error) {
	var ob types.OrderBook
	if err := json.Unmarshal(raw, &ob); err != nil {
		return nil, err
	}
	return ob, nil
}
