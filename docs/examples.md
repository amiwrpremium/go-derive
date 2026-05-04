# Examples

Two layers of examples — pick the one that fits your need:

## 1. godoc `Example_*` functions

Every public package has runnable examples that appear inline on
[pkg.go.dev](https://pkg.go.dev/github.com/amiwrpremium/go-derive). They
double as table-driven tests with `// Output:` blocks, so a regression
shows up in CI immediately.

86 example test files, 234 functions. One example file per source file:

```text
pkg/types/decimal_example_test.go         ↔  pkg/types/decimal.go
pkg/types/address_example_test.go         ↔  pkg/types/address.go
pkg/auth/local_signer_example_test.go     ↔  pkg/auth/local_signer.go
...
```

Run them all:

```bash
go test -count=1 -run='^Example' ./...
```

## 2. Runnable `examples/` programs

91 self-contained programs, one per directory, mirroring the SDK surface:

```text
examples/
├── derive/                4 programs   (facade)
├── auth/                  10 programs  (signing primitives)
├── rest/
│   ├── public/            7 programs
│   └── private/           20 programs
└── ws/
    ├── public/            12 programs (incl. 4 subscribe/)
    └── private/           27 programs (incl. 6 subscribe/)
```

Each `main.go` is ~10–25 lines because the env-loading boilerplate is
in [`examples/example/example.go`](../examples/example/example.go).

Run any one:

```bash
go run ./examples/rest/public/get_ticker
go run ./examples/ws/public/subscribe/orderbook

DERIVE_SESSION_KEY=0x... DERIVE_SUBACCOUNT=123 \
    go run ./examples/rest/private/get_open_orders
```

Build them all (CI does this on every PR):

```bash
go build ./examples/...
```

See [`examples/README.md`](../examples/README.md) for the full env-var
list.

## When to use which

| Use this | When |
|---|---|
| godoc `Example_*` | learning a single function's API in pkg.go.dev |
| `examples/` program | running it end-to-end against testnet, copy-pasting into your code |

The runnable programs assume Derive testnet by default; godoc examples
use only in-memory data so they pass offline in CI.
