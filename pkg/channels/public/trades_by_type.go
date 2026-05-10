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

// TradesByTypeTxStatus subscribes to all public trades for one
// (instrument_type, currency, tx_status) tuple — same shape as
// [TradesByType] but filtered server-side by on-chain transaction
// status.
//
// The dotted server-side channel name is:
//
//	trades.{instrument_type}.{currency}.{tx_status}
//
// Per the cockpit `channel_trades_instrument_type_currency_tx_status.rs`
// schema, only `settled` and `reverted` are documented filter values
// today — other [enums.TxStatus] values may be rejected by the engine.
// Same per-batch payload as [TradesByType]; pair with T = [[]types.Trade].
type TradesByTypeTxStatus struct {
	// InstrumentType narrows the stream to one product class.
	InstrumentType enums.InstrumentType
	// Currency is the underlying asset symbol (e.g. "BTC").
	Currency string
	// TxStatus is the on-chain transaction status to filter by.
	// Documented values: [enums.TxStatusSettled], [enums.TxStatusReverted].
	TxStatus enums.TxStatus
}

// Name returns the dotted server-side channel string.
func (t TradesByTypeTxStatus) Name() string {
	return fmt.Sprintf("trades.%s.%s.%s", t.InstrumentType, t.Currency, t.TxStatus)
}

// Decode parses an inbound notification payload into a [[]types.Trade].
// Same payload shape as the unfiltered [TradesByType] channel.
func (TradesByTypeTxStatus) Decode(raw json.RawMessage) (any, error) {
	var trades []types.Trade
	if err := json.Unmarshal(raw, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}
