# Numerics

Every numeric field in the SDK is `derive.Decimal` — never `float64`.

## Why not float64

Order sizes and prices on Derive go down to 18 decimal places. Marshalling
them through float64 truncates: `0.1 + 0.2 != 0.3` and that's the kind of
silent corruption that causes orders to be rejected for "invalid amount"
or, worse, accepted at slightly wrong prices.

`Decimal` is a thin wrapper over `shopspring/decimal`, the de-facto
high-precision decimal library. It JSON-encodes as a *string* — Derive's
wire convention — and the SDK's tests verify byte-for-byte round-trips.

## Constructors

```go
d, err := types.NewDecimal("65000.5")     // parse a string
d := types.MustDecimal("65000.5")         // panics on parse failure (constants/tests)
d := types.DecimalFromInt(42)             // exact int → decimal
```

## Arithmetic

`Decimal.Inner()` returns the underlying `shopspring/decimal.Decimal`:

```go
limit := tk.MarkPrice.Inner().Mul(decimal.RequireFromString("0.95"))
limitPrice, _ := types.NewDecimal(limit.String())
```

We deliberately don't expose arithmetic methods on `Decimal` directly —
shopspring's API is rich enough that wrapping it would just hide what's
available.

## ABI scaling for signing

Derive's on-chain modules use 18-decimal fixed-point: a price of `100`
goes on the wire as the integer `100 * 1e18`. The SDK's signing path
handles this via `internal/codec.DecimalToU256` and
`internal/codec.DecimalToI256`:

```go
price := decimal.RequireFromString("100.5")
scaled, err := codec.DecimalToU256(price)
// scaled = 100_500_000_000_000_000_000
```

The functions reject:

- **Negative input to `DecimalToU256`** (use `DecimalToI256` for signed).
- **Sub-1e-18 precision** — anything more precise than 18 dp can't be
  represented exactly at the contract's scale.

## On the wire

```json
{"price": "65000.5", "amount": "0.001"}
```

Decimals are JSON strings. The SDK's `UnmarshalJSON` also accepts JSON
numbers (e.g. `65000.5` without quotes) for resilience to upstream
schema drift, but the SDK always *emits* strings.

## Round-trip correctness

```go
d := types.MustDecimal("0.000000000000000001") // 1e-18
b, _ := json.Marshal(d)
// b == `"0.000000000000000001"`

var back derive.Decimal
_ = json.Unmarshal(b, &back)
// back.String() == "0.000000000000000001"  // identical
```

The fuzz test in `types.go/decimal_fuzz_test.go` runs this round-trip on
random input continuously.

## Common gotchas

- **`omitempty` doesn't omit a zero `Decimal`** — Go's `encoding/json`
  doesn't recognise custom-marshalled zero values as empty. A zero
  decimal serialises as the string `"0"` regardless of the tag. This is
  documented and asserted by `TestCandle_ZeroVolumeStillSerialized` in
  the test suite.
- **Don't compare `Decimal` with `==`** — use the underlying
  `Inner().Cmp(...)` because two equivalent decimals can have different
  internal representations (e.g. `1.0` vs `1` after JSON round-trip).
- **Negative amounts are valid for transfers, not for orders.** The
  signing helpers enforce this:
  `TradeModuleData.Hash` rejects negative `MaxFee`;
  `TransferModuleData.Hash` accepts a signed `Amount`.

## When to reach for raw `decimal.Decimal`

For arithmetic, comparison, or anything beyond round-trip storage, use
`d.Inner()` and work with the upstream library directly. The wrapper's
job is wire-format and zero-value safety; the math is shopspring's.
