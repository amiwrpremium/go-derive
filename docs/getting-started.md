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

    "github.com/amiwrpremium/go-derive"
    "github.com/amiwrpremium/go-derive"
)

func main() {
    c, err := derive.NewClient(derive.WithTestnet())
    if err != nil { log.Fatal(err) }
    defer c.Close()

    insts, err := c.REST.GetInstruments(context.Background(), "BTC", derive.InstrumentTypePerp)
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
import "github.com/amiwrpremium/go-derive"

signer, _ := derive.NewLocalSigner(os.Getenv("DERIVE_SESSION_KEY"))
c, _ := derive.NewClient(
    derive.WithTestnet(),
    derive.WithSigner(signer),
    derive.WithSubaccount(123),
)
```

## WebSocket subscription

```go
import (
    "github.com/amiwrpremium/go-derive"
    "github.com/amiwrpremium/go-derive"
    "github.com/amiwrpremium/go-derive"
)

c, _ := derive.NewClient(derive.WithTestnet())
defer c.Close()
c.WS.Connect(ctx)

sub, err := derive.Subscribe[derive.OrderBook](ctx, c.WS,
    derive.PublicOrderBook{Instrument: "BTC-PERP", Depth: 5})
defer sub.Close()

for ob := range sub.Updates() {
    fmt.Println(ob.Bids[0].Price)
}
```

## Environment variables

The SDK itself reads no env vars. Examples and integration tests do — see
[`examples/README.md`](../examples/README.md) and
[`test/README.md`](../test/README.md).

## Next steps

- [architecture.md](./architecture.md) for the layering rationale.
- [auth.md](./auth.md) for production signing setup.
- [`examples/`](../examples/) for 80 runnable programs covering every
  RPC method and channel.
