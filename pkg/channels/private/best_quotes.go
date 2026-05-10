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

// BestQuotes subscribes to the running stream of best-quote state
// for every open RFQ on one subaccount.
//
// The dotted server-side channel name (per
// docs.derive.xyz/reference/subaccount_id-best-quotes) is:
//
//	{subaccount_id}.best.quotes
//
// — note the unusual format: the subaccount id is the leading
// segment without a `subaccount.` prefix, and the suffix uses dots
// rather than the usual `_` (i.e. `best.quotes` not `best_quotes`).
//
// Each notification batch carries a [types.BestQuoteFeedEvent] per
// RFQ on the subaccount: either the engine's current best-quote
// state (Result) or the upstream RPC error (Error). Pair this
// descriptor with T = [[]types.BestQuoteFeedEvent] when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe].
type BestQuotes struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
}

// Name returns the dotted server-side channel string.
func (b BestQuotes) Name() string {
	return fmt.Sprintf("%d.best.quotes", b.SubaccountID)
}

// Decode parses an inbound notification payload into [[]types.BestQuoteFeedEvent].
func (BestQuotes) Decode(raw json.RawMessage) (any, error) {
	var events []types.BestQuoteFeedEvent
	if err := json.Unmarshal(raw, &events); err != nil {
		return nil, err
	}
	return events, nil
}
