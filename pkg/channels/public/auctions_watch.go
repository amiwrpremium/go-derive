// Package public declares Derive's public (no-auth) WebSocket subscription
// channels.
//
// This file holds the [AuctionsWatch] descriptor for the
// `auctions.watch` channel. See the package doc on [MarginWatch]
// for the cross-cutting interface contract.
package public

import (
	"encoding/json"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// AuctionsWatch subscribes to the platform-wide stream of liquidation
// auctions Derive is currently running.
//
// The dotted server-side channel name is the bare literal:
//
//	auctions.watch
//
// Per the docs (https://docs.derive.xyz/reference/auctions-watch),
// the channel takes no parameters — every subscribed client receives
// the same engine-wide stream. Consumers filter client-side on the
// [types.AuctionWatchEvent.SubaccountID] /
// [types.AuctionWatchEvent.Details.Currency] fields if they only
// care about a subset.
//
// Pair this descriptor with T = [types.AuctionWatchEvent] when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe]. Each
// notification carries one auction state transition.
type AuctionsWatch struct{}

// Name returns the dotted server-side channel string.
func (AuctionsWatch) Name() string { return "auctions.watch" }

// Decode parses an inbound notification payload into a [types.AuctionWatchEvent].
func (AuctionsWatch) Decode(raw json.RawMessage) (any, error) {
	var ev types.AuctionWatchEvent
	if err := json.Unmarshal(raw, &ev); err != nil {
		return nil, err
	}
	return ev, nil
}
