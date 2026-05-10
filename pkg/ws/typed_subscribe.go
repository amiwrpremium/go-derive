// Package ws is the WebSocket-backed client for Derive's JSON-RPC API.
//
// This file declares one Subscribe* method per documented Derive
// channel, each a thin convenience wrapper over the generic
// [Subscribe] function.
package ws

import (
	"context"

	"github.com/amiwrpremium/go-derive/pkg/channels/private"
	"github.com/amiwrpremium/go-derive/pkg/channels/public"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// The Subscribe* methods declared below are thin convenience wrappers
// over the generic [Subscribe] function. They exist so callers can subscribe in one
// line without hand-picking T and the matching descriptor:
//
//	// Generic form:
//	sub, _ := ws.Subscribe[[]types.Order](ctx, c, private.Orders{SubaccountID: 7})
//
//	// Typed form (this file):
//	sub, _ := c.SubscribeOrders(ctx, 7)
//
// The generic [Subscribe] remains the supported building block for
// custom or third-party channel descriptors; the typed methods here
// cover every channel listed at https://docs.derive.xyz/reference/.

// --------------------------------------------------------------------
// Public channels (no auth required).
// --------------------------------------------------------------------

// SubscribeMarginWatch streams the platform-wide stream of
// subaccounts whose maintenance margin has crossed the watch
// threshold. See [public.MarginWatch].
func (c *Client) SubscribeMarginWatch(ctx context.Context) (*Subscription[[]types.MarginWatch], error) {
	return Subscribe[[]types.MarginWatch](ctx, c, public.MarginWatch{})
}

// SubscribeAuctionsWatch streams the platform-wide state of ongoing
// liquidation auctions. See [public.AuctionsWatch].
func (c *Client) SubscribeAuctionsWatch(ctx context.Context) (*Subscription[types.AuctionWatchEvent], error) {
	return Subscribe[types.AuctionWatchEvent](ctx, c, public.AuctionsWatch{})
}

// SubscribeOrderBook streams incremental order-book updates for one
// instrument. Empty group defaults to "1" (no grouping); zero depth
// defaults to 10. See [public.OrderBook].
func (c *Client) SubscribeOrderBook(ctx context.Context, instrument, group string, depth int) (*Subscription[types.OrderBook], error) {
	return Subscribe[types.OrderBook](ctx, c, public.OrderBook{Instrument: instrument, Group: group, Depth: depth})
}

// SubscribeSpotFeed streams oracle-signed spot feeds for one
// currency. See [public.SpotFeed].
func (c *Client) SubscribeSpotFeed(ctx context.Context, currency string) (*Subscription[types.SpotFeed], error) {
	return Subscribe[types.SpotFeed](ctx, c, public.SpotFeed{Currency: currency})
}

// SubscribeTicker streams the full ticker payload for one
// instrument. Empty interval defaults to "1000". See [public.Ticker].
func (c *Client) SubscribeTicker(ctx context.Context, instrument, interval string) (*Subscription[types.InstrumentTickerFeed], error) {
	return Subscribe[types.InstrumentTickerFeed](ctx, c, public.Ticker{Instrument: instrument, Interval: interval})
}

// SubscribeTickerSlim streams the slim ticker payload for one
// instrument. Empty interval defaults to "1000".
// See [public.TickerSlim].
func (c *Client) SubscribeTickerSlim(ctx context.Context, instrument, interval string) (*Subscription[types.TickerSlim], error) {
	return Subscribe[types.TickerSlim](ctx, c, public.TickerSlim{Instrument: instrument, Interval: interval})
}

// SubscribeTrades streams public trade prints on one instrument.
// See [public.Trades].
func (c *Client) SubscribeTrades(ctx context.Context, instrument string) (*Subscription[[]types.Trade], error) {
	return Subscribe[[]types.Trade](ctx, c, public.Trades{Instrument: instrument})
}

// SubscribeTradesByType streams public trade prints aggregated by
// instrument type and currency. See [public.TradesByType].
func (c *Client) SubscribeTradesByType(ctx context.Context, instrumentType enums.InstrumentType, currency string) (*Subscription[[]types.Trade], error) {
	return Subscribe[[]types.Trade](ctx, c, public.TradesByType{InstrumentType: instrumentType, Currency: currency})
}

// SubscribeTradesByTypeWithStatus is like SubscribeTradesByType but
// also filters by on-chain transaction status.
// See [public.TradesByTypeTxStatus].
func (c *Client) SubscribeTradesByTypeWithStatus(ctx context.Context, instrumentType enums.InstrumentType, currency string, txStatus enums.TxStatus) (*Subscription[[]types.Trade], error) {
	return Subscribe[[]types.Trade](ctx, c, public.TradesByTypeTxStatus{InstrumentType: instrumentType, Currency: currency, TxStatus: txStatus})
}

// --------------------------------------------------------------------
// Private channels (require [Client.Login]).
// --------------------------------------------------------------------

// SubscribeBalances streams balance updates for one subaccount.
// See [private.Balances].
func (c *Client) SubscribeBalances(ctx context.Context, subaccountID int64) (*Subscription[types.Balance], error) {
	return Subscribe[types.Balance](ctx, c, private.Balances{SubaccountID: subaccountID})
}

// SubscribeOrders streams order lifecycle events for one
// subaccount. See [private.Orders].
func (c *Client) SubscribeOrders(ctx context.Context, subaccountID int64) (*Subscription[[]types.Order], error) {
	return Subscribe[[]types.Order](ctx, c, private.Orders{SubaccountID: subaccountID})
}

// SubscribeBestQuotes streams the running best-quote state for every
// open RFQ on one subaccount. See [private.BestQuotes].
func (c *Client) SubscribeBestQuotes(ctx context.Context, subaccountID int64) (*Subscription[[]types.BestQuoteFeedEvent], error) {
	return Subscribe[[]types.BestQuoteFeedEvent](ctx, c, private.BestQuotes{SubaccountID: subaccountID})
}

// SubscribeRFQs streams RFQ lifecycle events for one wallet across
// every subaccount it owns. See [private.RFQs].
func (c *Client) SubscribeRFQs(ctx context.Context, wallet string) (*Subscription[[]types.RFQ], error) {
	return Subscribe[[]types.RFQ](ctx, c, private.RFQs{Wallet: wallet})
}

// SubscribeQuotes streams quote events for one subaccount.
// See [private.Quotes].
func (c *Client) SubscribeQuotes(ctx context.Context, subaccountID int64) (*Subscription[[]types.Quote], error) {
	return Subscribe[[]types.Quote](ctx, c, private.Quotes{SubaccountID: subaccountID})
}

// SubscribeSubaccountTrades streams trade events for one
// subaccount. See [private.Trades].
func (c *Client) SubscribeSubaccountTrades(ctx context.Context, subaccountID int64) (*Subscription[[]types.Trade], error) {
	return Subscribe[[]types.Trade](ctx, c, private.Trades{SubaccountID: subaccountID})
}

// SubscribeSubaccountTradesByStatus is like SubscribeSubaccountTrades
// but also filters by on-chain transaction status.
// See [private.TradesByTxStatus].
func (c *Client) SubscribeSubaccountTradesByStatus(ctx context.Context, subaccountID int64, txStatus enums.TxStatus) (*Subscription[[]types.Trade], error) {
	return Subscribe[[]types.Trade](ctx, c, private.TradesByTxStatus{SubaccountID: subaccountID, TxStatus: txStatus})
}
