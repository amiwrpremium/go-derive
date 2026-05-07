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

// TickerSlim subscribes to ticker updates for one instrument.
//
// The dotted server-side channel name is:
//
//	ticker_slim.{instrument}.{interval}
//
// Interval is the update cadence in milliseconds — `100` or `1000`. Empty
// defaults to `1000`. (The legacy `ticker.{instrument}.{interval}ms` pattern
// was deprecated by Derive in favour of this one, which delivers a more
// compact wire payload — single-letter field names like `b`/`B` for the
// best-bid price/amount.)
//
// Pair this descriptor with T = [types.TickerSlim] when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe].
type TickerSlim struct {
	// Instrument is the market name (e.g. "BTC-PERP").
	Instrument string
	// Interval is the update cadence in milliseconds as a string.
	// Allowed values: `100`, `1000`. Empty means `1000`.
	Interval string
}

// Name returns the dotted server-side channel string.
func (t TickerSlim) Name() string {
	i := t.Interval
	if i == "" {
		i = "1000"
	}
	return fmt.Sprintf("ticker_slim.%s.%s", t.Instrument, i)
}

// Decode parses an inbound notification payload into a [types.TickerSlim].
func (TickerSlim) Decode(raw json.RawMessage) (any, error) {
	var t types.TickerSlim
	if err := json.Unmarshal(raw, &t); err != nil {
		return nil, err
	}
	return t, nil
}
