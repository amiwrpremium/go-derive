// Package ws is the WebSocket-backed client for Derive's JSON-RPC API.
//
// This file declares one Subscribe* method per documented Derive
// channel, each a thin convenience wrapper over the generic
// [Subscribe] function. They exist so callers can subscribe in one
// line without hand-building the dotted channel name.
package ws

import (
	"context"
	"fmt"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// The Subscribe* methods declared below cover every channel listed
// at https://docs.derive.xyz/reference/. The generic [Subscribe]
// remains available for custom or yet-undocumented channels.
//
// All methods accept variadic [SubscribeOption] arguments to tune
// buffer size, drop policy, and the error handler — see the
// pkg/ws/subscribe_options.go declarations for the trade-offs.

// --------------------------------------------------------------------
// Public channels (no auth required).
// --------------------------------------------------------------------

// SubscribeMarginWatch streams the platform-wide stream of
// subaccounts whose maintenance margin has crossed the watch
// threshold. Wire channel: `margin.watch`.
func (c *Client) SubscribeMarginWatch(ctx context.Context, opts ...SubscribeOption) (*Subscription[[]types.MarginWatch], error) {
	return Subscribe(ctx, c, "margin.watch", decodeJSON[[]types.MarginWatch], opts...)
}

// SubscribeAuctionsWatch streams the platform-wide state of ongoing
// liquidation auctions. Wire channel: `auctions.watch`.
func (c *Client) SubscribeAuctionsWatch(ctx context.Context, opts ...SubscribeOption) (*Subscription[types.AuctionWatchEvent], error) {
	return Subscribe(ctx, c, "auctions.watch", decodeJSON[types.AuctionWatchEvent], opts...)
}

// SubscribeOrderBook streams incremental order-book updates for one
// instrument. Empty group defaults to "1" (no grouping); zero depth
// defaults to 10. Wire channel:
// `orderbook.{instrument}.{group}.{depth}`.
func (c *Client) SubscribeOrderBook(ctx context.Context, instrument, group string, depth int, opts ...SubscribeOption) (*Subscription[types.OrderBook], error) {
	if group == "" {
		group = "1"
	}
	if depth == 0 {
		depth = 10
	}
	return Subscribe(ctx, c,
		fmt.Sprintf("orderbook.%s.%s.%d", instrument, group, depth),
		decodeJSON[types.OrderBook], opts...)
}

// SubscribeSpotFeed streams oracle-signed spot feeds for one
// currency. Wire channel: `spot_feed.{currency}`.
func (c *Client) SubscribeSpotFeed(ctx context.Context, currency string, opts ...SubscribeOption) (*Subscription[types.SpotFeed], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("spot_feed.%s", currency),
		decodeJSON[types.SpotFeed], opts...)
}

// SubscribeTicker streams the full ticker payload for one
// instrument. Empty interval defaults to "1000". Wire channel:
// `ticker.{instrument}.{interval}`.
func (c *Client) SubscribeTicker(ctx context.Context, instrument, interval string, opts ...SubscribeOption) (*Subscription[types.InstrumentTickerFeed], error) {
	if interval == "" {
		interval = "1000"
	}
	return Subscribe(ctx, c,
		fmt.Sprintf("ticker.%s.%s", instrument, interval),
		decodeJSON[types.InstrumentTickerFeed], opts...)
}

// SubscribeTickerSlim streams the slim ticker payload for one
// instrument. Empty interval defaults to "1000". Wire channel:
// `ticker_slim.{instrument}.{interval}`.
func (c *Client) SubscribeTickerSlim(ctx context.Context, instrument, interval string, opts ...SubscribeOption) (*Subscription[types.TickerSlim], error) {
	if interval == "" {
		interval = "1000"
	}
	return Subscribe(ctx, c,
		fmt.Sprintf("ticker_slim.%s.%s", instrument, interval),
		decodeJSON[types.TickerSlim], opts...)
}

// SubscribeTrades streams public trade prints on one instrument.
// Wire channel: `trades.{instrument}`.
func (c *Client) SubscribeTrades(ctx context.Context, instrument string, opts ...SubscribeOption) (*Subscription[[]types.Trade], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("trades.%s", instrument),
		decodeJSON[[]types.Trade], opts...)
}

// SubscribeTradesByType streams public trade prints aggregated by
// instrument type and currency. Wire channel:
// `trades.{instrument_type}.{currency}`.
func (c *Client) SubscribeTradesByType(ctx context.Context, instrumentType enums.InstrumentType, currency string, opts ...SubscribeOption) (*Subscription[[]types.Trade], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("trades.%s.%s", instrumentType, currency),
		decodeJSON[[]types.Trade], opts...)
}

// SubscribeTradesByTypeWithStatus is like [Client.SubscribeTradesByType]
// but also filters by on-chain transaction status. Wire channel:
// `trades.{instrument_type}.{currency}.{tx_status}`.
func (c *Client) SubscribeTradesByTypeWithStatus(ctx context.Context, instrumentType enums.InstrumentType, currency string, txStatus enums.TxStatus, opts ...SubscribeOption) (*Subscription[[]types.Trade], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("trades.%s.%s.%s", instrumentType, currency, txStatus),
		decodeJSON[[]types.Trade], opts...)
}

// --------------------------------------------------------------------
// Private channels (require [Client.Login]).
// --------------------------------------------------------------------

// SubscribeBalances streams balance updates for one subaccount.
// Wire channel: `{subaccount_id}.balances`.
func (c *Client) SubscribeBalances(ctx context.Context, subaccountID int64, opts ...SubscribeOption) (*Subscription[types.Balance], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("%d.balances", subaccountID),
		decodeJSON[types.Balance], opts...)
}

// SubscribeOrders streams order lifecycle events for one
// subaccount. Wire channel: `{subaccount_id}.orders`.
func (c *Client) SubscribeOrders(ctx context.Context, subaccountID int64, opts ...SubscribeOption) (*Subscription[[]types.Order], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("%d.orders", subaccountID),
		decodeJSON[[]types.Order], opts...)
}

// SubscribeBestQuotes streams the running best-quote state for every
// open RFQ on one subaccount. Wire channel:
// `{subaccount_id}.best.quotes`.
func (c *Client) SubscribeBestQuotes(ctx context.Context, subaccountID int64, opts ...SubscribeOption) (*Subscription[[]types.BestQuoteFeedEvent], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("%d.best.quotes", subaccountID),
		decodeJSON[[]types.BestQuoteFeedEvent], opts...)
}

// SubscribeRFQs streams RFQ lifecycle events for one wallet across
// every subaccount it owns. Wire channel: `{wallet}.rfqs`.
func (c *Client) SubscribeRFQs(ctx context.Context, wallet string, opts ...SubscribeOption) (*Subscription[[]types.RFQ], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("%s.rfqs", wallet),
		decodeJSON[[]types.RFQ], opts...)
}

// SubscribeQuotes streams quote events for one subaccount. Wire
// channel: `{subaccount_id}.quotes`.
func (c *Client) SubscribeQuotes(ctx context.Context, subaccountID int64, opts ...SubscribeOption) (*Subscription[[]types.Quote], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("%d.quotes", subaccountID),
		decodeJSON[[]types.Quote], opts...)
}

// SubscribeSubaccountTrades streams trade events for one
// subaccount. Wire channel: `{subaccount_id}.trades`.
func (c *Client) SubscribeSubaccountTrades(ctx context.Context, subaccountID int64, opts ...SubscribeOption) (*Subscription[[]types.Trade], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("%d.trades", subaccountID),
		decodeJSON[[]types.Trade], opts...)
}

// SubscribeSubaccountTradesByStatus is like
// [Client.SubscribeSubaccountTrades] but also filters by on-chain
// transaction status. Wire channel:
// `{subaccount_id}.trades.{tx_status}`.
func (c *Client) SubscribeSubaccountTradesByStatus(ctx context.Context, subaccountID int64, txStatus enums.TxStatus, opts ...SubscribeOption) (*Subscription[[]types.Trade], error) {
	return Subscribe(ctx, c,
		fmt.Sprintf("%d.trades.%s", subaccountID, txStatus),
		decodeJSON[[]types.Trade], opts...)
}
