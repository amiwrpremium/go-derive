// Package public.
package public

import (
	"encoding/json"
	"fmt"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// TradesByType subscribes to all public trades for one (instrument_type,
// currency) combination — e.g. every perp print on BTC, every option on ETH.
//
// The dotted server-side channel name is:
//
//	trades.{instrument_type}.{currency}
//
// Where InstrumentType is one of [enums.InstrumentType] (perp, option, erc20)
// and Currency is the underlying symbol (BTC, ETH, …).
//
// Pair this descriptor with T = [[]types.Trade]. Each notification carries
// a batch of trades that printed in the same window across every instrument
// matching the (type, currency) tuple — useful for index-level analytics
// without subscribing per-instrument.
type TradesByType struct {
	// InstrumentType narrows the stream to one product class.
	InstrumentType enums.InstrumentType
	// Currency is the underlying asset symbol (e.g. "BTC").
	Currency string
}

// Name returns the dotted server-side channel string.
func (t TradesByType) Name() string {
	return fmt.Sprintf("trades.%s.%s", t.InstrumentType, t.Currency)
}

// Decode parses an inbound notification payload into a [[]types.Trade].
func (TradesByType) Decode(raw json.RawMessage) (any, error) {
	var trades []types.Trade
	if err := json.Unmarshal(raw, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}
