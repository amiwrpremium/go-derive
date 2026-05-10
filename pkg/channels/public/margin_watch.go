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

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// MarginWatch subscribes to the platform-wide stream of subaccounts
// whose maintenance margin has crossed the watch threshold.
//
// The dotted server-side channel name is the bare literal:
//
//	margin.watch
//
// Per the docs (https://docs.derive.xyz/reference/margin-watch),
// the channel takes no parameters — every subscribed client receives
// the same engine-wide stream. Consumers filter client-side on the
// [types.MarginWatch.MarginType] / [types.MarginWatch.SubaccountID]
// fields if they only care about a subset.
//
// Pair this descriptor with T = [[]types.MarginWatch] when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe]. Each
// notification carries a batch of at-risk subaccounts captured at
// the same `valuation_timestamp`.
type MarginWatch struct{}

// Name returns the dotted server-side channel string.
func (MarginWatch) Name() string { return "margin.watch" }

// Decode parses an inbound notification payload into [[]types.MarginWatch].
func (MarginWatch) Decode(raw json.RawMessage) (any, error) {
	var events []types.MarginWatch
	if err := json.Unmarshal(raw, &events); err != nil {
		return nil, err
	}
	return events, nil
}
