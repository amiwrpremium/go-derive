# Subscriptions

The WebSocket transport supports pub/sub channels. The SDK gives you
typed access to every documented channel via the descriptors in
`pkg/channels/public` and `pkg/channels/private`.

## The pattern

```go
sub, err := ws.Subscribe[T](ctx, c, descriptor)
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
| `public.OrderBook{Instrument, Group, Depth}` | `orderbook.{i}.{g}.{d}` | `types.OrderBook` |
| `public.Trades{Instrument}` | `trades.{i}` | `[]types.Trade` |
| `public.TradesByType{InstrumentType, Currency}` | `trades.{type}.{currency}` | `[]types.Trade` |
| `public.TickerSlim{Instrument, Interval}` | `ticker_slim.{i}.{interval}` | `types.TickerSlim` |
| `public.SpotFeed{Currency}` | `spot_feed.{currency}` | `types.SpotFeed` |

`OrderBook` accepts `Group ∈ {1, 10, 100}` (price-bucket size) and
`Depth ∈ {1, 10}` (levels per side). `TickerSlim` accepts
`Interval ∈ {100, 1000}` (milliseconds; default 1000). Derive removed
the legacy `instruments.{c}.{t}` and `ticker.{i}.{interval}ms` channels —
poll `public/get_instruments` over REST/WS-RPC for the former, use
`TickerSlim` for the latter.

Private:

| Descriptor | Channel name | T |
|---|---|---|
| `private.Orders{SubaccountID}` | `subaccount.{id}.orders` | `[]types.Order` |
| `private.Balances{SubaccountID}` | `subaccount.{id}.balances` | `types.Balance` |
| `private.Trades{SubaccountID}` | `subaccount.{id}.trades` | `[]types.Trade` |
| `private.RFQs{Wallet}` | `wallet.{address}.rfqs` | `[]types.RFQ` |
| `private.Quotes{SubaccountID}` | `subaccount.{id}.quotes` | `[]types.Quote` |

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
err := ws.SubscribeFunc(ctx, c, public.OrderBook{Instrument: "BTC-PERP"},
    func(ob types.OrderBook) {
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
documented and asserted in `pkg/ws/subscribe.go`). For a reliable queue,
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
