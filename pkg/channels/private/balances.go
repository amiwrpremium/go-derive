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

// Balances subscribes to collateral and total-equity updates on one
// subaccount.
//
// The dotted server-side channel name (per
// https://docs.derive.xyz/reference/subaccount_id-balances) is:
//
//	{subaccount_id}.balances
//
// Pair this descriptor with T = [types.Balance] (a single struct, not a
// slice) when calling [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe].
type Balances struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
}

// Name returns the dotted server-side channel string.
func (b Balances) Name() string { return fmt.Sprintf("%d.balances", b.SubaccountID) }

// Decode parses an inbound notification payload into a [types.Balance].
func (Balances) Decode(raw json.RawMessage) (any, error) {
	var bal types.Balance
	if err := json.Unmarshal(raw, &bal); err != nil {
		return nil, err
	}
	return bal, nil
}
