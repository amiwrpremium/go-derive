# Subscriptions

The WebSocket transport supports pub/sub channels. The SDK gives you
typed access to every documented channel via the descriptors in
`channels.go` and `channels.go`.

## The pattern

```go
sub, err := derive.Subscribe[T](ctx, c, descriptor)
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

1. **Generic `T`** — the descriptor decodes JSON into a typed Go value;
   the type parameter on `Subscribe` removes the cast at the call site.
2. **`select` against `ctx.Done()`** — caller controls cancellation.
3. **`<-chan T` close** — terminal signal. Always check `sub.Err()` after
   the channel closes.

## Channel descriptors

Every descriptor implements `channels.Channel`:

```go
type Channel interface {
    Name() string
    Decode(raw json.RawMessage) (any, error)
}
```

Public:

| Descriptor | Channel name | T |
|---|---|---|
| `derive.PublicOrderBook{Instrument, Group, Depth}` | `orderbook.{i}.{g}.{d}` | `derive.OrderBook` |
| `derive.PublicTrades{Instrument}` | `trades.{i}` | `[]derive.Trade` |
| `derive.PublicTradesByType{InstrumentType, Currency}` | `trades.{type}.{currency}` | `[]derive.Trade` |
| `derive.PublicTickerSlim{Instrument, Interval}` | `ticker_slim.{i}.{interval}` | `derive.TickerSlim` |
| `derive.PublicSpotFeed{Currency}` | `spot_feed.{currency}` | `derive.SpotFeed` |

`OrderBook` accepts `Group ∈ {1, 10, 100}` (price-bucket size) and
`Depth ∈ {1, 10}` (levels per side). `TickerSlim` accepts
`Interval ∈ {100, 1000}` (milliseconds; default 1000). Derive removed
the legacy `instruments.{c}.{t}` and `ticker.{i}.{interval}ms` channels —
poll `public/get_instruments` over REST/WS-RPC for the former, use
`TickerSlim` for the latter.

Private:

| Descriptor | Channel name | T |
|---|---|---|
| `derive.PrivateOrders{SubaccountID}` | `subaccount.{id}.orders` | `[]derive.Order` |
| `derive.PrivateBalances{SubaccountID}` | `subaccount.{id}.balances` | `derive.Balance` |
| `derive.PrivateTrades{SubaccountID}` | `subaccount.{id}.trades` | `[]derive.Trade` |
| `derive.PrivateRFQs{Wallet}` | `wallet.{address}.rfqs` | `[]derive.RFQ` |
| `derive.PrivateQuotes{SubaccountID}` | `subaccount.{id}.quotes` | `[]derive.Quote` |

Note: there is **no** `subaccount.{id}.positions` channel. Poll
`private/get_positions` or derive position state from the trades feed.
RFQs are wallet-scoped (one stream per signer address), not
subaccount-scoped.

Private channels require `c.WS.Login(ctx)` first.

## Callback variant: `SubscribeFunc`

When channel-receive is awkward (e.g. integrating with an existing event
loop), `SubscribeFunc` drives a callback synchronously and returns when
the context cancels:

```go
err := derive.SubscribeFunc(ctx, c, derive.PublicOrderBook{Instrument: "BTC-PERP"},
    func(ob derive.OrderBook) {
        process(ob)
    })
// err is ctx.Err() or the terminal subscription error
```

The callback runs synchronously, so back-pressure on the caller is
back-pressure on the subscription. That's intentional — `Subscribe[T]`'s
buffered channel drops oldest events under back-pressure (best-effort
fan-out, not a reliable queue).

## Reconnect

When `WithReconnect(true)` (the default), the WS transport re-dials on
drops, re-runs the `OnReconnect` hook (`Login` for authenticated clients),
then re-issues every active `subscribe` so user-facing `Subscription[T]`
channels stay open across the gap.

See [reconnection.md](./reconnection.md).

## Buffer sizing

Each `Subscription[T]` has a 256-event buffer. If you process events
slower than they arrive, newer events are dropped (drop-oldest is *not*
the policy — current implementation drops new on a full buffer; this is
documented and asserted in `ws.go/subscribe.go`). For a reliable queue,
use `SubscribeFunc` and apply your own bounded queueing inside the
callback.

## Subscribing twice to the same channel

Calling `Subscribe[T]` twice with the same channel name on the same
client returns the *same* underlying subscription — no extra RPC is
issued. This is a feature: it lets independent components in the same
process share a stream cheaply.

## Unsubscribing

`sub.Close()` issues an unsubscribe RPC best-effort and drains the typed
channel. Idempotent — calling it twice is harmless.

## See also

- One runnable program per descriptor under
  [`examples/ws/{public,private}/subscribe/<channel>/`](../examples/ws/).
- A multi-channel demux pattern in
  [`examples/ws/public/subscribe/multi/`](../examples/ws/public/subscribe/multi/).
- An auto-reconnect-resilience demo in
  [`examples/ws/public/subscribe/reconnect/`](../examples/ws/public/subscribe/reconnect/).
- A private multi-channel demo (orders + positions + balances) in
  [`examples/ws/private/subscribe/multi/`](../examples/ws/private/subscribe/multi/).
