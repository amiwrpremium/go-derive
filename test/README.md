# Integration tests

These tests hit Derive's live network. They're gated by the `integration`
build tag, so the default `go test ./...` is unaffected.

## Quickstart

```bash
# Public-only subset — no creds needed.
go test -tags=integration -count=1 -run='^TestPublic_|^TestWS_PublicConnect|^TestWS_OrderBookSubscribe|^TestWS_TickerSubscribe|^TestCross_' ./test/...

# All tests except live order placement.
DERIVE_SESSION_KEY=0x... DERIVE_SUBACCOUNT=123 \
  go test -tags=integration -count=1 ./test/...

# Add live order placement (testnet only — never run against mainnet).
DERIVE_SESSION_KEY=0x... DERIVE_SUBACCOUNT=123 DERIVE_BASE_ASSET=0x... DERIVE_RUN_LIVE_ORDERS=1 \
  go test -tags=integration -count=1 -run='^TestPrivate_PlaceAndCancel|^TestPrivate_OrderEventsArrive' ./test/...
```

Or use the Makefile shortcut:

```bash
make test-integration
```

## Environment variables

| Var | Required for | Purpose |
|---|---|---|
| `DERIVE_NETWORK` | all (default `testnet`) | `mainnet` or `testnet` |
| `DERIVE_INSTRUMENT` | all (default `BTC-PERP`) | instrument used by the public + private tests |
| `DERIVE_SESSION_KEY` | private tests | hex-encoded session-key private key (with or without `0x`) |
| `DERIVE_OWNER` | private tests when distinct from session-key address | smart-account owner address |
| `DERIVE_SUBACCOUNT` | private tests | numeric subaccount id |
| `DERIVE_BASE_ASSET` | live-order tests | on-chain asset address used for trade-module signing |
| `DERIVE_RUN_LIVE_ORDERS` | live-order tests | set to `1` to opt in to real order placement |

Tests that require credentials call `t.Skip(...)` when the relevant env
vars are missing — running without creds simply runs the public subset.

## Safety rails

- Default network is **testnet**; the test code never assumes mainnet.
- Live order placement is **double-gated**: requires both private creds
  and `DERIVE_RUN_LIVE_ORDERS=1`. Tests place a buy 5% below mark (won't
  fill) and cancel it.
- Per-test timeouts are 30 seconds (10s for snapshot subscriptions, 60s
  for the order-flow test that needs round-trip).
- The order-flow tests re-cancel via `defer` so a panic still leaves the
  account clean.

## Test layout

| File | Scope |
|---|---|
| `integration_test.go` | env loading + client builders + skip helpers |
| `rest_public_integration_test.go` | get_time, get_currencies, get_instruments, get_ticker, get_orderbook, get_trade_history |
| `rest_private_integration_test.go` | get_subaccount(s), get_open_orders, get_positions, get_collateral, paginated history endpoints |
| `ws_public_integration_test.go` | Connect, OrderBook/Ticker/Trades subscriptions |
| `ws_private_integration_test.go` | Login, subaccount Orders/Balances/Positions subscriptions |
| `place_order_integration_test.go` | Place + cancel via REST and WS; verify the order appears on the WS Orders channel |
| `cross_transport_integration_test.go` | REST vs WS equivalence checks; `derive.Client` facade smoke test |

## CI

Integration tests are not run on push or PR (they depend on Derive testnet
availability and configured secrets). They run on manual `workflow_dispatch`
only — see `.github/workflows/integration.yml`.
