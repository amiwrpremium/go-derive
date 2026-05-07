// Package private — see balances.go for the overview.
package private

import (
	"encoding/json"
	"fmt"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// Orders subscribes to order lifecycle events on one subaccount.
//
// The dotted server-side channel name is:
//
//	subaccount.{id}.orders
//
// Pair this descriptor with T = [[]types.Order] when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe].
type Orders struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
}

// Name returns the dotted server-side channel string.
func (o Orders) Name() string { return fmt.Sprintf("subaccount.%d.orders", o.SubaccountID) }

// Decode parses an inbound notification payload into a [[]types.Order].
func (Orders) Decode(raw json.RawMessage) (any, error) {
	var orders []types.Order
	if err := json.Unmarshal(raw, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}
