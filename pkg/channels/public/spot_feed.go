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

// SpotFeed subscribes to oracle price-feed updates for one currency.
//
// The dotted server-side channel name is:
//
//	spot_feed.{currency}
//
// Currency is the underlying asset symbol (e.g. "BTC", "ETH"). The
// payload includes the current oracle price plus the 24-hour-prior
// reading, so consumers can render a delta without an extra round trip.
//
// Pair this descriptor with T = [derive.SpotFeed].
type SpotFeed struct {
	// Currency is the underlying asset symbol (e.g. "BTC").
	Currency string
}

// Name returns the dotted server-side channel string.
func (s SpotFeed) Name() string { return fmt.Sprintf("spot_feed.%s", s.Currency) }

// Decode parses an inbound notification payload into a [derive.SpotFeed].
func (SpotFeed) Decode(raw json.RawMessage) (any, error) {
	var sf derive.SpotFeed
	if err := json.Unmarshal(raw, &sf); err != nil {
		return nil, err
	}
	return sf, nil
}
