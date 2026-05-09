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

	"github.com/amiwrpremium/go-derive"
)

// RFQs subscribes to lifecycle updates for RFQs initiated by one wallet.
//
// The dotted server-side channel name is:
//
//	wallet.{address}.rfqs
//
// RFQs on Derive are wallet-scoped — a single signer address sees every
// RFQ it issued across all of its subaccounts. Address must be a 0x-prefixed
// 20-byte hex string in standard EIP-55 form.
//
// Pair with T = [[]derive.RFQ].
type RFQs struct {
	// Wallet is the owner address as a 0x-prefixed hex string.
	Wallet string
}

// Name returns the dotted server-side channel string.
func (r RFQs) Name() string { return fmt.Sprintf("wallet.%s.rfqs", r.Wallet) }

// Decode parses an inbound notification payload into a [[]derive.RFQ].
func (RFQs) Decode(raw json.RawMessage) (any, error) {
	var rfqs []derive.RFQ
	if err := json.Unmarshal(raw, &rfqs); err != nil {
		return nil, err
	}
	return rfqs, nil
}

// Quotes subscribes to quote updates received against the subaccount's
// outstanding RFQs.
//
// The dotted server-side channel name is:
//
//	subaccount.{id}.quotes
//
// Pair with T = [[]derive.Quote].
type Quotes struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
}

// Name returns the dotted server-side channel string.
func (q Quotes) Name() string { return fmt.Sprintf("subaccount.%d.quotes", q.SubaccountID) }

// Decode parses an inbound notification payload into a [[]derive.Quote].
func (Quotes) Decode(raw json.RawMessage) (any, error) {
	var quotes []derive.Quote
	if err := json.Unmarshal(raw, &quotes); err != nil {
		return nil, err
	}
	return quotes, nil
}
