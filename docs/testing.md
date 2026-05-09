# Testing

Four overlapping test layers. All run on every PR; the integration suite
is opt-in.

## 1. Unit tests

One test file per source file. ~100 `_test.go` files.

```bash
go test -count=1 ./...
go test -race -count=1 ./...      # also catches concurrency bugs
```

Coverage gates:

| Package | Floor |
|---|---|
| `auth.go` | ≥ 95% |
| `errors.go` | 100% |
| `types.go` | ≥ 98% |
| `internal/jsonrpc` | 100% |
| `internal/codec` | 100% |
| Total | ≥ 88% |

Coverage report:

```bash
go test -coverprofile=coverage.txt -covermode=atomic -coverpkg=./...,./internal/... ./...
go tool cover -func=coverage.txt | tail -1
```

Codecov receives the per-PR diff via [ci.yml](../.github/workflows/ci.yml).

## 2. godoc Example_* tests

234 functions across 86 files, one example file per source file. Each
example carries an `// Output:` block so a behaviour regression makes
`go test` fail.

```bash
go test -count=1 -run='^Example' ./...
```

## 3. Fuzz tests

Native Go `Fuzz*` tests guard the parsers most likely to crash on
adversarial input:

| Package | Targets |
|---|---|
| `types.go` | `FuzzNewDecimal`, `FuzzDecimal_UnmarshalJSON`, `FuzzNewAddress`, `FuzzNewTxHash`, `FuzzMillisTime_UnmarshalJSON`, `FuzzOrderBookLevel_UnmarshalJSON` |
| `auth.go` | `FuzzNewLocalSigner` |
| `errors.go` | `FuzzAPIError_UnmarshalJSON` |
| `internal/jsonrpc` | `FuzzIsNotification`, `FuzzDecodeResult` |

The seed corpus runs as part of every `go test`. To fuzz indefinitely
(out of CI):

```bash
go test -fuzz=FuzzIsNotification -fuzztime=10m ./internal/jsonrpc
```

## 4. Integration tests (opt-in)

Live testnet calls under [`test/`](../test/), gated by the `integration`
build tag.

```bash
# Public-only subset (no creds)
make test-integration

# All except live order placement
DERIVE_SESSION_KEY=0x... DERIVE_SUBACCOUNT=123 \
    go test -tags=integration -count=1 ./test/...

# Live order placement (testnet only)
DERIVE_SESSION_KEY=0x... DERIVE_SUBACCOUNT=123 \
DERIVE_BASE_ASSET=0x... DERIVE_RUN_LIVE_ORDERS=1 \
    go test -tags=integration -count=1 -run='^TestPrivate_PlaceAndCancel' ./test/...
```

Safety rails:

- Default network is testnet. Tests never auto-promote to mainnet.
- Live orders require both private creds *and*
  `DERIVE_RUN_LIVE_ORDERS=1`. Tests place 5% below mark (won't fill) and
  cancel.
- Per-test 30s timeouts, 60s for the order-flow round-trip.

See [`test/README.md`](../test/README.md).

## Test fakes

| Fake | Used by | What it does |
|---|---|---|
| `internal/testutil.FakeTransport` | every `methods.go/*_test.go` | implements `transport.Transport`; programmable per-method handlers |
| `internal/testutil.MockServer` | `rest.go` | httptest-based JSON-RPC HTTP mock |
| `internal/testutil.MockWSServer` | `ws.go`, `internal/transport/ws_test.go` | `gorilla/websocket.Upgrader`-based JSON-RPC mock with notification injection |

## Verifying every source has a test

```bash
comm -23 \
  <(find pkg internal -name '*.go' ! -name '*_test.go' ! -name 'doc.go' \
     ! -name 'state.go' ! -name 'util.go' ! -name 'channel.go' ! -name 'transport.go' \
     ! -path '*/testutil/*' | sort) \
  <(find pkg internal -name '*_test.go' ! -name 'internal_helpers_test.go' ! -name 'internal_time_test.go' | sed 's/_test\.go/.go/' | sort)
# expected: empty
```

## Useful commands

```bash
make test           # vet + go test ./...
make test-race      # add -race
make cover          # coverage summary
make test-integration   # integration suite (testnet)
```
