// Package channels declares the typed WebSocket subscription channels
// Derive supports.
//
// Each channel descriptor (see pkg/channels/public and pkg/channels/private)
// implements [Channel]: it knows the dotted server-side name to subscribe
// to and how to decode an inbound notification payload into a typed Go
// value. The descriptors are passed into
// [Subscribe] and
// [SubscribeFunc].
//
// # Public vs private
//
// Public channels (pkg/channels/public) need no authentication. Private
// channels (pkg/channels/private) need a logged-in WebSocket — the
// SubaccountID field on each private descriptor scopes the stream to
// one subaccount.
package derive

import (
	"encoding/json"
	"fmt"
)

// Channel is the contract a subscription descriptor must satisfy.
//
// Implementations live in pkg/channels/public and pkg/channels/private.
// Third-party implementations are welcome — anything that produces a
// dotted name and decodes JSON to a typed value will work.
type Channel interface {
	// Name returns the dot-separated server-side channel name, e.g.
	// "trades.BTC-PERP" or "subaccount.123.orders".
	Name() string

	// Decode turns a raw notification payload into a typed value. The
	// concrete return type is descriptor-specific; pass the matching T to
	// ws.Subscribe to consume it without an explicit type assertion.
	Decode(raw json.RawMessage) (any, error)
}

// OrderBook subscribes to incremental order-book updates for one instrument.
//
// The dotted server-side channel name is:
//
//	orderbook.{instrument}.{group}.{depth}
//
// where Group is the price grouping (e.g. "10" for nearest-10-ticks; "1"
// for no grouping) and Depth is the number of levels returned per side.
// Empty Group defaults to "1"; zero Depth defaults to 10.
//
// Pair this descriptor with T = [OrderBook] when calling
// [Subscribe].
type PublicOrderBook struct {
	// Instrument is the market name (e.g. "BTC-PERP").
	Instrument string
	// Group is the price-grouping string. Empty means "1" (no grouping).
	Group string
	// Depth is the number of book levels per side. Zero means 10.
	Depth int
}

// Name returns the dotted server-side channel string. See [PublicOrderBook] for
// the format.
func (o PublicOrderBook) Name() string {
	g := o.Group
	if g == "" {
		g = "1"
	}
	d := o.Depth
	if d == 0 {
		d = 10
	}
	return fmt.Sprintf("orderbook.%s.%s.%d", o.Instrument, g, d)
}

// Decode parses one inbound notification payload into a [OrderBook].
func (PublicOrderBook) Decode(raw json.RawMessage) (any, error) {
	var ob OrderBook
	if err := json.Unmarshal(raw, &ob); err != nil {
		return nil, err
	}
	return ob, nil
}

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
// Pair this descriptor with T = [SpotFeed].
type PublicSpotFeed struct {
	// Currency is the underlying asset symbol (e.g. "BTC").
	Currency string
}

// Name returns the dotted server-side channel string.
func (s PublicSpotFeed) Name() string { return fmt.Sprintf("spot_feed.%s", s.Currency) }

// Decode parses an inbound notification payload into a [SpotFeed].
func (PublicSpotFeed) Decode(raw json.RawMessage) (any, error) {
	var sf SpotFeed
	if err := json.Unmarshal(raw, &sf); err != nil {
		return nil, err
	}
	return sf, nil
}

// TickerSlim subscribes to ticker updates for one instrument.
//
// The dotted server-side channel name is:
//
//	ticker_slim.{instrument}.{interval}
//
// Interval is the update cadence in milliseconds — `100` or `1000`. Empty
// defaults to `1000`. (The legacy `ticker.{instrument}.{interval}ms` pattern
// was deprecated by Derive in favour of this one, which delivers a more
// compact wire payload — single-letter field names like `b`/`B` for the
// best-bid price/amount.)
//
// Pair this descriptor with T = [TickerSlim] when calling
// [Subscribe].
type PublicTickerSlim struct {
	// Instrument is the market name (e.g. "BTC-PERP").
	Instrument string
	// Interval is the update cadence in milliseconds as a string.
	// Allowed values: `100`, `1000`. Empty means `1000`.
	Interval string
}

// Name returns the dotted server-side channel string.
func (t PublicTickerSlim) Name() string {
	i := t.Interval
	if i == "" {
		i = "1000"
	}
	return fmt.Sprintf("ticker_slim.%s.%s", t.Instrument, i)
}

// Decode parses an inbound notification payload into a [TickerSlim].
func (PublicTickerSlim) Decode(raw json.RawMessage) (any, error) {
	var t TickerSlim
	if err := json.Unmarshal(raw, &t); err != nil {
		return nil, err
	}
	return t, nil
}

// TradesByType subscribes to all public trades for one (instrument_type,
// currency) combination — e.g. every perp print on BTC, every option on ETH.
//
// The dotted server-side channel name is:
//
//	trades.{instrument_type}.{currency}
//
// Where InstrumentType is one of [InstrumentType] (perp, option, erc20)
// and Currency is the underlying symbol (BTC, ETH, …).
//
// Pair this descriptor with T = [[]Trade]. Each notification carries
// a batch of trades that printed in the same window across every instrument
// matching the (type, currency) tuple — useful for index-level analytics
// without subscribing per-instrument.
type PublicTradesByType struct {
	// InstrumentType narrows the stream to one product class.
	InstrumentType InstrumentType
	// Currency is the underlying asset symbol (e.g. "BTC").
	Currency string
}

// Name returns the dotted server-side channel string.
func (t PublicTradesByType) Name() string {
	return fmt.Sprintf("trades.%s.%s", t.InstrumentType, t.Currency)
}

// Decode parses an inbound notification payload into a [[]Trade].
func (PublicTradesByType) Decode(raw json.RawMessage) (any, error) {
	var trades []Trade
	if err := json.Unmarshal(raw, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}

// Trades subscribes to public trade prints on one instrument.
//
// The dotted server-side channel name is:
//
//	trades.{instrument}
//
// Pair this descriptor with T = []Trade when calling
// [Subscribe]. Each notification
// carries a batch of [Trade] events that printed in the same window.
type PublicTrades struct {
	// Instrument is the market name (e.g. "BTC-PERP").
	Instrument string
}

// Name returns the dotted server-side channel string.
func (t PublicTrades) Name() string { return fmt.Sprintf("trades.%s", t.Instrument) }

// Decode parses an inbound notification payload into a [[]Trade].
func (PublicTrades) Decode(raw json.RawMessage) (any, error) {
	var trades []Trade
	if err := json.Unmarshal(raw, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}

// Balances subscribes to collateral and total-equity updates on one
// subaccount.
//
// The dotted server-side channel name is:
//
//	subaccount.{id}.balances
//
// Pair this descriptor with T = [Balance] (a single struct, not a
// slice) when calling [Subscribe].
type PrivateBalances struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
}

// Name returns the dotted server-side channel string.
func (b PrivateBalances) Name() string { return fmt.Sprintf("subaccount.%d.balances", b.SubaccountID) }

// Decode parses an inbound notification payload into a [Balance].
func (PrivateBalances) Decode(raw json.RawMessage) (any, error) {
	var bal Balance
	if err := json.Unmarshal(raw, &bal); err != nil {
		return nil, err
	}
	return bal, nil
}

// Orders subscribes to order lifecycle events on one subaccount.
//
// The dotted server-side channel name is:
//
//	subaccount.{id}.orders
//
// Pair this descriptor with T = [[]Order] when calling
// [Subscribe].
type PrivateOrders struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
}

// Name returns the dotted server-side channel string.
func (o PrivateOrders) Name() string { return fmt.Sprintf("subaccount.%d.orders", o.SubaccountID) }

// Decode parses an inbound notification payload into a [[]Order].
func (PrivateOrders) Decode(raw json.RawMessage) (any, error) {
	var orders []Order
	if err := json.Unmarshal(raw, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

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
// Pair with T = [[]RFQ].
type PrivateRFQs struct {
	// Wallet is the owner address as a 0x-prefixed hex string.
	Wallet string
}

// Name returns the dotted server-side channel string.
func (r PrivateRFQs) Name() string { return fmt.Sprintf("wallet.%s.rfqs", r.Wallet) }

// Decode parses an inbound notification payload into a [[]RFQ].
func (PrivateRFQs) Decode(raw json.RawMessage) (any, error) {
	var rfqs []RFQ
	if err := json.Unmarshal(raw, &rfqs); err != nil {
		return nil, err
	}
	return rfqs, nil
}

// Quotes subscribes to quote updates received against the subaccount's
// outstanding PrivateRFQs.
//
// The dotted server-side channel name is:
//
//	subaccount.{id}.quotes
//
// Pair with T = [[]Quote].
type PrivateQuotes struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
}

// Name returns the dotted server-side channel string.
func (q PrivateQuotes) Name() string { return fmt.Sprintf("subaccount.%d.quotes", q.SubaccountID) }

// Decode parses an inbound notification payload into a [[]Quote].
func (PrivateQuotes) Decode(raw json.RawMessage) (any, error) {
	var quotes []Quote
	if err := json.Unmarshal(raw, &quotes); err != nil {
		return nil, err
	}
	return quotes, nil
}

// Trades subscribes to fill events for one subaccount.
//
// The dotted server-side channel name is:
//
//	subaccount.{id}.trades
//
// Pair this descriptor with T = [[]Trade] when calling
// [Subscribe]. Each notification
// is a batch of fills since the last update.
type PrivateTrades struct {
	// SubaccountID scopes the stream to one subaccount.
	SubaccountID int64
}

// Name returns the dotted server-side channel string.
func (t PrivateTrades) Name() string { return fmt.Sprintf("subaccount.%d.trades", t.SubaccountID) }

// Decode parses an inbound notification payload into a [[]Trade].
func (PrivateTrades) Decode(raw json.RawMessage) (any, error) {
	var trades []Trade
	if err := json.Unmarshal(raw, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}
