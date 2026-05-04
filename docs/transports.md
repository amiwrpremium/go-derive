# Transports: REST vs WebSocket

Derive's API is JSON-RPC 2.0 over **both HTTP and WebSocket**. The SDK
exposes both as `pkg/rest.Client` and `pkg/ws.Client`. They share the
underlying method definitions (`internal/methods.API`), so identical Go
code compiles against either.

## When to pick which

| Need | Use |
|---|---|
| One-shot read (instruments, ticker) | REST |
| Streaming order book / trades / ticker | WebSocket (subscriptions) |
| Order placement at low latency | WebSocket (no per-call TLS handshake) |
| Long-running maker process | WebSocket — keep one socket alive |
| Throwaway batch script | REST |

REST is simpler to reason about; WebSocket is faster and the only way to
stream live data.

## Order-book retrieval — REST has no endpoint

Derive removed the REST `public/get_orderbook` method. Two replacements,
depending on what you need:

- **Top-of-book + 5 % depth (snapshot, REST-friendly)**: call
  `GetTicker` — `best_ask_price`/`best_ask_amount`,
  `best_bid_price`/`best_bid_amount`, and the
  `five_percent_{bid,ask}_depth` fields cover most workflows that
  previously hit the REST orderbook.
- **Full L2 (any depth, streaming or one-shot)**: subscribe to the
  WebSocket `orderbook.<instrument>.<group>.<depth>` channel via
  `pkg/channels/public.OrderBook`. See
  [`examples/ws/public/subscribe/orderbook/`](../examples/ws/public/subscribe/orderbook/).
  Cancel after the first message if you only want a snapshot.

## Both at once: `pkg/derive.Client`

```go
c, _ := derive.NewClient(derive.WithTestnet(), derive.WithSigner(s), derive.WithSubaccount(1))
defer c.Close()

// REST for setup
insts, _ := c.REST.GetInstruments(ctx, "BTC", enums.InstrumentTypePerp)

// WS for streaming
c.WS.Connect(ctx)
c.WS.Login(ctx)
sub, _ := ws.Subscribe[types.OrderBook](ctx, c.WS, public.OrderBook{Instrument: "BTC-PERP"})
```

## Cross-transport equivalence

Because both clients call into the same `internal/methods.API`, identical
parameters produce identical wire calls. The integration suite has
explicit cross-transport tests:

- `TestCross_GetInstruments` — REST and WS return the same instrument set.
- `TestCross_GetTicker` — marks within 1% of each other (small lag tolerance).

See [`test/cross_transport_integration_test.go`](../test/cross_transport_integration_test.go).

## Authentication shape

REST puts the EIP-191 timestamp signature in headers
(`X-LyraWallet`/`X-LyraTimestamp`/`X-LyraSignature`) on every request.
WebSocket does it once via the `public/login` RPC after `Connect`. The
SDK's auto-reconnect path re-runs `Login` on every reconnect so private
subscriptions resume automatically.

## Method coverage

Every method in `internal/methods.API` is reachable from both clients.
Subscription channels are WebSocket-only — REST has no streaming primitive.

## Rate limiting

Each client owns its own token-bucket limiter (default: 10 TPS, burst 5×).
REST and WS limit independently — that's a design choice; if you saturate
one transport, the other still has headroom. See
[rate-limiting.md](./rate-limiting.md) to tune.
