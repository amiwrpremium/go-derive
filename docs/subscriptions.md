# Subscriptions

The WebSocket transport supports pub/sub channels. The SDK gives you
typed access to every documented channel via the `Subscribe*` methods
on `*ws.Client` — one per documented channel.

## The pattern

```go
sub, err := c.WS.SubscribeOrderBook(ctx, "BTC-PERP", ws.GroupDefault, ws.DepthDefault)
if err != nil { return err }
defer sub.Close()

for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case event, ok := <-sub.Updates():
        if !ok {
            return sub.Err()
        }
        process(event)
    }
}
```

Three things going on:

1. **Typed `Subscribe*` method** — the method's return type is
   `*Subscription[T]` for the channel's documented payload `T`. No
   cast at the call site.
2. **`select` against `ctx.Done()`** — caller controls cancellation.
3. **`<-chan T` close** — terminal signal. Always check `sub.Err()`
   after the channel closes.

## Public channels

| Method | Wire channel | `T` |
|---|---|---|
| `SubscribeOrderBook(instrument, group, depth)` | `orderbook.{instrument}.{group}.{depth}` | `types.OrderBook` |
| `SubscribeSpotFeed(currency)` | `spot_feed.{currency}` | `types.SpotFeed` |
| `SubscribeTickerSlim(instrument, interval)` | `ticker_slim.{instrument}.{interval}` | `types.TickerSlim` |
| `SubscribeTicker(instrument, interval)` ⚠ | `ticker.{instrument}.{interval}` | `types.InstrumentTickerFeed` |
| `SubscribeTrades(instrument)` | `trades.{instrument}` | `[]types.Trade` |
| `SubscribeTradesByType(instrumentType, currency)` | `trades.{instrument_type}.{currency}` | `[]types.Trade` |
| `SubscribeTradesByTypeWithStatus(instrumentType, currency, txStatus)` | `trades.{instrument_type}.{currency}.{tx_status}` | `[]types.Trade` |
| `SubscribeMarginWatch()` | `margin.watch` | `[]types.MarginWatch` |
| `SubscribeAuctionsWatch()` | `auctions.watch` | `types.AuctionWatchEvent` |

⚠ `SubscribeTicker` is deprecated since 2025-12-01 (upstream Derive
notice). Use `SubscribeTickerSlim` for new code. The verbose-ticker
method still works; `staticcheck` flags callers via SA1019 to surface
the deprecation at compile time.

`SubscribeOrderBook` accepts `group ∈ {"1", "10", "100"}` (price-bucket
size) and `depth ∈ {1, 10, 20, 100}` (levels per side). Pass `""` / `0`
to get `GroupDefault` (`"1"`) and `DepthDefault` (`10`). The
`Interval100` / `Interval1000` constants in `pkg/ws` cover the ticker
intervals.

## Private channels

| Method | Wire channel | `T` |
|---|---|---|
| `SubscribeOrders(subaccountID)` | `{subaccount_id}.orders` | `[]types.Order` |
| `SubscribeBalances(subaccountID)` | `{subaccount_id}.balances` | `[]types.BalanceUpdate` |
| `SubscribeSubaccountTrades(subaccountID)` | `{subaccount_id}.trades` | `[]types.Trade` |
| `SubscribeSubaccountTradesByStatus(subaccountID, txStatus)` | `{subaccount_id}.trades.{tx_status}` | `[]types.Trade` |
| `SubscribeQuotes(subaccountID)` | `{subaccount_id}.quotes` | `[]types.Quote` |
| `SubscribeBestQuotes(subaccountID)` | `{subaccount_id}.best.quotes` | `[]types.BestQuoteFeedEvent` |
| `SubscribeRFQs(wallet)` | `{wallet}.rfqs` | `[]types.RFQ` |

`SubscribeBalances` decodes into `[]types.BalanceUpdate` (per-row delta
events), not a full snapshot — the channel is event-driven.

There is **no** `{subaccount_id}.positions` channel. Poll
`private/get_positions` over REST/WS-RPC, or derive position state
from `SubscribeSubaccountTrades`. RFQs are wallet-scoped (one stream
per signer address), not subaccount-scoped.

Private channels require `c.WS.Login(ctx)` first.

## Generic `Subscribe[T]` for custom or undocumented channels

If you need to subscribe to a channel without a typed wrapper (e.g. a
new channel template added upstream before the SDK catches up), the
underlying generic is exported:

```go
sub, err := ws.Subscribe[types.OrderBook](ctx, c.WS,
    "orderbook.BTC-PERP.1.10",
    ws.DecodeJSON[types.OrderBook])
```

Caller supplies the channel name and a decoder (`ws.DecodeJSON[T]` is
the standard JSON decoder; you can plug in your own for binary
payloads or alternate schemas).

## Callback variant: `SubscribeFunc`

When channel-receive is awkward (e.g. integrating with an existing
event loop), `SubscribeFunc` drives a callback synchronously and
returns when the context cancels or the subscription terminates:

```go
err := ws.SubscribeFunc[types.OrderBook](ctx, c.WS,
    "orderbook.BTC-PERP.1.10",
    ws.DecodeJSON[types.OrderBook],
    func(ob types.OrderBook) {
        process(ob)
    })
// err is ctx.Err() or the terminal subscription error
```

The callback runs synchronously on the subscription's read pump, so
back-pressure on the caller is back-pressure on the subscription.
That's intentional — `Subscribe[T]`'s buffered channel can drop
events under back-pressure (see "Buffer sizing").

## Reconnect

When `WithReconnect(true)` (the default), the WS transport re-dials on
drops, re-runs the post-dial hook (`Login` for authenticated clients),
then re-issues every active `subscribe` so user-facing
`Subscription[T]` channels stay open across the gap.

See [reconnection.md](./reconnection.md).

## Buffer sizing

Each `Subscription[T]` has a default 256-event buffer. Tune it via
`ws.WithBufferSize(n)`. If the consumer falls behind:

- Default `DropPolicy` is configurable via `ws.WithDropPolicy(...)`.
  See `pkg/ws/subscribe_options.go` for the exact semantics.
- For a reliable queue, use `SubscribeFunc` and apply your own
  bounded queueing inside the callback — the synchronous callback
  applies natural back-pressure to the subscription.

## Subscribing twice to the same channel

Calling any `Subscribe*` method twice with the same channel name on
the same client returns the *same* underlying subscription — no
extra RPC is issued. This is a feature: it lets independent
components in the same process share a stream cheaply.

## Unsubscribing

`sub.Close()` issues an unsubscribe RPC best-effort and drains the
typed channel. Idempotent — calling it twice is harmless.

## See also

- A runnable program per typed `Subscribe*` method lives under
  [`examples/ws/{public,private}/subscribe/<channel>/`](../examples/ws/).
- A multi-channel demux pattern in
  [`examples/ws/public/subscribe/multi/`](../examples/ws/public/subscribe/multi/).
- An auto-reconnect-resilience demo in
  [`examples/ws/public/subscribe/reconnect/`](../examples/ws/public/subscribe/reconnect/).
