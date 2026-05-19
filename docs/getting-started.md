# Getting started

## Install

```bash
go get github.com/amiwrpremium/go-derive
```

Requires Go 1.25+.

## First program (public, no credentials)

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/amiwrpremium/go-derive/pkg/derive"
    "github.com/amiwrpremium/go-derive/pkg/enums"
)

func main() {
    c, err := derive.NewClient(derive.WithTestnet())
    if err != nil { log.Fatal(err) }
    defer c.Close()

    insts, err := c.REST.GetInstruments(context.Background(), "BTC", enums.InstrumentTypePerp)
    if err != nil { log.Fatal(err) }
    fmt.Println(len(insts), "BTC perps")
}
```

That's it — no auth needed for market-data calls.

## With credentials (private endpoints)

Derive uses **session keys**: a hot key registered on-chain by the
smart-account owner. For development you can use the same key for both
("LocalSigner"); production deployments should use `SessionKeySigner` so
the long-lived owner key never lives in the trading process.

```go
import "github.com/amiwrpremium/go-derive/pkg/auth"

signer, _ := auth.NewLocalSigner(os.Getenv("DERIVE_SESSION_KEY"))
c, _ := derive.NewClient(
    derive.WithTestnet(),
    derive.WithSigner(signer),
    derive.WithSubaccount(123),
)
```

## WebSocket subscription

```go
import (
    "github.com/amiwrpremium/go-derive/pkg/derive"
    "github.com/amiwrpremium/go-derive/pkg/ws"
)

ctx := context.Background()

c, _ := derive.NewClient(derive.WithTestnet(), derive.WithConnectWS(true))
defer c.Close()

sub, err := c.WS.SubscribeOrderBook(ctx, "BTC-PERP", ws.GroupDefault, ws.DepthDefault)
if err != nil { log.Fatal(err) }
defer sub.Close()

for ob := range sub.Updates() {
    if len(ob.Bids) > 0 {
        fmt.Println(ob.Bids[0].Price)
    }
}
```

`SubscribeOrderBook` is one of ~16 typed `Subscribe*` methods on
`*ws.Client` — one per documented channel. See
[subscriptions.md](./subscriptions.md) for the full list and the
non-default values for `group`/`depth`/`interval`. The generic
`ws.Subscribe[T](ctx, c, channelName, decoder, opts...)` form is also
available for channels not covered by a typed wrapper.

## Environment variables

The SDK itself reads no env vars. Examples and integration tests do — see
[`examples/README.md`](../examples/README.md) and
[`test/README.md`](../test/README.md).

## Next steps

- [architecture.md](./architecture.md) for the layering rationale.
- [auth.md](./auth.md) for production signing setup.
- [`examples/`](../examples/) for runnable programs covering every
  RPC method and channel.
