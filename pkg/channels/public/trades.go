// Package public — see orderbook.go for the overview.
package public

import (
	"encoding/json"
	"fmt"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// Trades subscribes to public trade prints on one instrument.
//
// The dotted server-side channel name is:
//
//	trades.{instrument}
//
// Pair this descriptor with T = []types.Trade when calling
// [github.com/amiwrpremium/go-derive/pkg/ws.Subscribe]. Each notification
// carries a batch of [types.Trade] events that printed in the same window.
type Trades struct {
	// Instrument is the market name (e.g. "BTC-PERP").
	Instrument string
}

// Name returns the dotted server-side channel string.
func (t Trades) Name() string { return fmt.Sprintf("trades.%s", t.Instrument) }

// Decode parses an inbound notification payload into a [[]types.Trade].
func (Trades) Decode(raw json.RawMessage) (any, error) {
	var trades []types.Trade
	if err := json.Unmarshal(raw, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}
