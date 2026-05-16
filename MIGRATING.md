# Migration guide

This document tracks every breaking change made on the 0.x line. Each
section spells out the before/after shape so callers can search for the
identifier they were using and find what to switch to.

For the release-by-release record (commit messages, dates, PR links)
see [CHANGELOG.md](./CHANGELOG.md). This guide is the *human* version:
what to change in your code.

---

## From any 0.x to 1.0 — consolidated TL;DR

Grouped by category so you can fix all of one kind at once. Each
entry links to the per-release section below for the full before/after.

### Method signatures

- **Every method that takes one or more non-context arguments now takes a single typed struct** (`*Input` for mutations, `*Query` for reads). Paginated reads take the query struct first and `types.PageRequest` second. → [v0.11.0](#v0110), [v0.15.0](#v0150), [v0.19.0](#v0190).
- **Read methods return values, not pointers.** `(T, error)` everywhere except for live resources (Client, Subscription, Signer). → [v0.10.0](#v0100).
- **`PlaceOrder`/`PlaceAlgoOrder`/`PlaceTriggerOrder` return `(Order, []Trade, error)`** instead of `(Order, error)` — the matched fills come back with the order. → [v0.14.0](#v0140).
- **`CancelOrder` returns the cancelled `types.Order`** instead of just `error`. → [v0.14.0](#v0140).

### Type renames

- **`signer.Address()` → `signer.SessionAddress()`**, **`signer.Owner()` → `signer.OwnerAddress()`**. → [v0.17.0](#v0170).
- **`types.GetOrdersFilter` → `types.OrdersQuery`** (value, not pointer). → [v0.19.0](#v0190).
- **`OrderParams` / `CancelOrderParams` / `ReplaceOrderParams` removed** — use `PlaceOrderInput` / `CancelOrderInput` / `ReplaceOrderInput`. → [v0.21.0](#v0210).
- **`GetCollateral` → `GetCollaterals`**, **`OrderQuote` → `GetOrderQuote`**, **`OrderQuotePublic` → `GetPublicOrderQuote`**. → [v0.22.0](#v0220).
- **`pkg/contracts` is gone.** Use `go-ethereum` directly for on-chain operations (deposit, withdraw, session-key registration). → [v0.20.0](#v0200).

### WebSocket / subscriptions

- **Cancelling the ctx passed to `Subscribe` / `SubscribeFunc` / `SubscribeInto` now terminates the subscription.** Pass `context.Background()` to keep a subscription alive after the call site's ctx expires. → [v0.9.0](#v090).
- **Channel wire names match the docs.** `margin.watch` (was `margin_watch`), `{id}.balances` (was `subaccount.{id}.balances`), and 5 other similar drifts. Only affects callers using raw channel strings; the typed `c.Subscribe*` methods always emitted the right names. → [v0.8.x](#v08x).
- **`SubscribeBalances` delivers `[]types.BalanceUpdate`** (not `types.Balance`). The channel is event-based, not a snapshot. → [v0.14.0](#v0140).

### Engine-side payload

- **RFQ `SendQuote` / `ExecuteQuote` / `ReplaceQuote` payloads are signed inside the SDK.** The four signature fields come off the input DTOs; the SDK fills them. → [v0.13.0](#v0130).
- **`SendRFQ`** takes `types.SendRFQInput` instead of positional args, and its wire key is now `max_total_cost` (was `max_total_fee`, which the engine silently ignored). → [v0.14.0](#v0140).

### Constructor vocabulary

- **`NewMillisTime(time.Time)` → `MillisTimeFromTime(time.Time)`**, and `NewMillisTime` now parses strings. New `MillisTimeFromMillis(int64)` and `MustMillisTime(string)`. → [v0.15.0](#v0150).
- **New `AddressFromCommon`, `DecimalFromShopspring`, `MustTxHash`, `TxHashFromCommon`** constructors. Additive — old constructors still exist. → [v0.15.0](#v0150).

---

## v0.22.0

**[#149](https://github.com/amiwrpremium/go-derive/pull/149)** — method-name normalization.

```go
// before                       // after
c.GetCollateral(ctx)            c.GetCollaterals(ctx)
c.OrderQuote(ctx, in)           c.GetOrderQuote(ctx, in)
c.OrderQuotePublic(ctx, in)     c.GetPublicOrderQuote(ctx, in)
```

The examples directories `examples/{rest,ws}/private/get_collateral/` are also renamed to `get_collaterals/`.

---

## v0.21.0

**[#147](https://github.com/amiwrpremium/go-derive/pull/147)** — unused `OrderParams` builder pattern removed.

```go
// before
op := types.NewOrderParams("BTC-PERP", buy, limit, amount, price).
    WithLabel("alpha").WithReduceOnly()
co := types.NewCancelOrderParams(subID).WithOrderID("O-1")
rp := types.NewReplaceOrderParams("O-1", op)

// after
op := types.PlaceOrderInput{
    InstrumentName: "BTC-PERP", Direction: buy, OrderType: limit,
    Amount: amount, LimitPrice: price, Label: "alpha", ReduceOnly: true,
}
co := types.CancelOrderInput{OrderID: "O-1"}
rp := types.ReplaceOrderInput{OrderIDToCancel: "O-1", PlaceOrderInput: op}
```

`types.ErrInvalidParams` remains importable; it moved to `pkg/types/validate.go`.

---

## v0.20.0

**[#145](https://github.com/amiwrpremium/go-derive/pull/145)** — `pkg/contracts` deleted.

The package shipped three interfaces (`Depositor`, `Withdrawer`, `SessionKeyManager`) whose every method returned `ErrNotImplemented`. No replacement in this SDK — for on-chain deposits, withdrawals, and session-key registration, use `go-ethereum` (or any EVM toolchain) directly against the Derive contracts. Once funds are on-chain, every order / RFQ / quote / cancel flow runs through this SDK.

---

## v0.19.0

**[#143](https://github.com/amiwrpremium/go-derive/pull/143)** — `GetOrders` / `GetOrderHistory` query+page shape.

```go
// before
c.GetOrders(ctx, types.PageRequest{PageSize: 10}, &types.GetOrdersFilter{InstrumentName: "BTC-PERP"})
c.GetOrderHistory(ctx, types.PageRequest{PageSize: 10}, types.OrderHistoryQuery{...})

// after
c.GetOrders(ctx, types.OrdersQuery{InstrumentName: "BTC-PERP"}, types.PageRequest{PageSize: 10})
c.GetOrderHistory(ctx, types.OrderHistoryQuery{...}, types.PageRequest{PageSize: 10})
```

Type rename `GetOrdersFilter` → `OrdersQuery`, swap of arg order, and the filter is now a value (no `&`).

---

## v0.18.0

**[#140](https://github.com/amiwrpremium/go-derive/pull/140)** — WS subscription ergonomics. Mostly additive; one rename worth noting:

No method signatures changed in this release. New additions: named depth/group/interval constants in `pkg/ws`, `WithOnReconnect` callback, validation of zero subaccount IDs on private `Subscribe*` methods, improved docs on multi-subscription drop policy.

---

## v0.17.0

**[#138](https://github.com/amiwrpremium/go-derive/pull/138)** — `Signer` interface rename.

```go
// before                       // after
signer.Address()                signer.SessionAddress()
signer.Owner()                  signer.OwnerAddress()
```

For `LocalSigner` both return the same address (the EOA). For `SessionKeySigner` they diverge — `SessionAddress` is the session key's EOA, `OwnerAddress` is the smart-account wallet that delegated to it.

---

## v0.16.0

**[#136](https://github.com/amiwrpremium/go-derive/pull/136)** — additive. New `rest.WithHTTPTimeout(d time.Duration)` option. No migration required.

---

## v0.15.0

Two breaking changes in this release.

### [#130](https://github.com/amiwrpremium/go-derive/pull/130) — Strict signature consolidation

Every method on `*methods.API` that takes one or more non-context arguments now takes a single typed `*Input` or `*Query` struct. 44 methods converted. The rule:

| Shape | When |
|---|---|
| `(ctx)` | no params |
| `(ctx, types.*Input{...})` | mutates state |
| `(ctx, types.*Query{...})` | reads state |
| `(ctx, types.*Query{...}, page)` | paginated reads — page stays separate |

Exceptions stay positional: variadic strings (`PreloadInstruments`), and the non-network helpers (`InvalidateInstrumentCache`, `SetTradeModule`, `SetRFQModule`).

Examples:

```go
// before                                                  // after
c.GetTicker(ctx, "BTC-PERP")                               c.GetTicker(ctx, types.TickerQuery{Name: "BTC-PERP"})
c.GetInstruments(ctx, "BTC", enums.InstrumentTypePerp)     c.GetInstruments(ctx, types.InstrumentsQuery{Currency: "BTC", Kind: enums.InstrumentTypePerp})
c.CancelOrder(ctx, "BTC-PERP", "O-12345")                  c.CancelOrder(ctx, types.CancelOrderInput{InstrumentName: "BTC-PERP", OrderID: "O-12345"})
c.MarginWatch(ctx, subID, false, false)                    c.MarginWatch(ctx, types.MarginWatchQuery{SubaccountID: subID})
```

See the PR for the full 44-method conversion table.

### [#133](https://github.com/amiwrpremium/go-derive/pull/133) — Symmetric identifier constructors

Unified the `New*` / `Must*` / `*FromX` vocabulary across `Address`, `Decimal`, `TxHash`, `MillisTime`.

The breaking part is `MillisTime`:

```go
// before                                  // after
types.NewMillisTime(t)                     types.MillisTimeFromTime(t)
types.NewMillisTime(time.UnixMilli(ms))    types.MillisTimeFromMillis(ms)
// NewMillisTime(string) didn't exist      types.NewMillisTime("1700000000000")
//                                         types.MustMillisTime("1700000000000")
```

New constructors (additive): `AddressFromCommon`, `DecimalFromShopspring`, `MustTxHash`, `TxHashFromCommon`.

---

## v0.14.1 — additive only

**[#126](https://github.com/amiwrpremium/go-derive/pull/126)** + **[#127](https://github.com/amiwrpremium/go-derive/pull/127)** — instrument cache, auto-pagination `*All` companions, error categorization helpers. No migration required.

---

## v0.14.0

Four breaking changes bundled to clear the docs-vs-SDK audit.

### [#119](https://github.com/amiwrpremium/go-derive/pull/119) — WS alignment

- `BalanceUpdate` struct field names match the documented wire shape: `name`, `new_balance`, `previous_balance`, `update_type` (was `asset_name`, `amount`, `previous_amount`, plus extras the docs never described).
- `SubscribeBalances` returns `*Subscription[[]types.BalanceUpdate]` (was `*Subscription[types.Balance]`). The channel publishes arrays of per-asset delta events.
- `CancelOrder` returns `(types.Order, error)` (was `error`). The cancelled order's `cancel_reason` and final timestamp come back with it.
- `OrderBook` gains a `PublishID` field for gap detection.

### [#122](https://github.com/amiwrpremium/go-derive/pull/122) — DTO field completeness

~50 documented fields filled in across `Ticker`, `Instrument`, `Position`, `Collateral`, `InstrumentTickerSlim`, `Trade`, plus the per-kind `Perp` / `Option` / `ERC20` detail blocks. All `omitempty` so existing payloads still decode.

One field rename: **`PerpDetails.AggregateFundingRate`** JSON tag changed from `aggregate_funding_rate` to `aggregate_funding`. Decoded payloads will now actually populate this field.

### [#123](https://github.com/amiwrpremium/go-derive/pull/123) — Optional read filters

```go
// before                                       // after
c.GetStatistics(ctx, "BTC-PERP")                c.GetStatistics(ctx, types.StatisticsQuery{InstrumentName: "BTC-PERP"})
c.GetPublicTradeHistory(ctx, page, "BTC-PERP")  c.GetPublicTradeHistory(ctx, types.PublicTradeHistoryQuery{InstrumentName: "BTC-PERP"}, page)
c.GetTradeHistory(ctx, page, fromTs, toTs)      c.GetTradeHistory(ctx, types.TradeHistoryQuery{FromTimestamp: fromTs, ToTimestamp: toTs}, page)
```

Each query type has the full filter set the docs describe (`Currency`, `EndTime`, `InstrumentType`, `Wallet`, `TxStatus`, etc.).

### [#125](https://github.com/amiwrpremium/go-derive/pull/125) — Optional signed-action fields + `SendRFQ` wire fix

```go
// before
order, err := c.PlaceOrder(ctx, types.PlaceOrderInput{...})

// after — same call, but you also get the matched fills
order, trades, err := c.PlaceOrder(ctx, types.PlaceOrderInput{...})
```

Same return-shape change for `PlaceAlgoOrder` and `PlaceTriggerOrder`. `PlaceOrderInput` gains `Client`, `IsAtomicSigning`, `ReferralCode`, `RejectPostOnly`, `RejectTimestamp`, `ExtraFee`.

```go
// before
c.SendRFQ(ctx, legs, maxFee)

// after — typed input, more documented fields, AND a wire-key bug fix
c.SendRFQ(ctx, types.SendRFQInput{
    Legs:         legs,
    MaxTotalCost: cost, // was max_total_fee on the wire — engine silently ignored it
    // optional: Counterparties, PreferredDirection, ReducingDirection, Label, MinTotalCost, ...
})
```

**Watch out** for the wire-key change: existing callers that passed a cost cap had it silently ignored. After upgrading, the cap is enforced server-side, so RFQs that previously matched at any cost may now be rejected.

---

## v0.13.0

**[#117](https://github.com/amiwrpremium/go-derive/pull/117)** — SDK signs RFQ payloads internally.

`SendQuoteInput`, `ExecuteQuoteInput`, `ReplaceQuoteInput` lose the four signature fields (`Nonce`, `Signature`, `Signer`, `SignatureExpirySec`). The SDK populates them using the configured signer, exactly like `PlaceOrder`.

```go
// before — caller pre-signed
in := types.SendQuoteInput{...}
in.Nonce = nonce
in.Signature = "0x..."
in.Signer = signerAddr
in.SignatureExpirySec = expiry
c.SendQuote(ctx, in)

// after
c.SendQuote(ctx, types.SendQuoteInput{ /* no signature fields */ })
```

---

## v0.12.0

**[#108](https://github.com/amiwrpremium/go-derive/pull/108)** — `examples/example` helper package removed.

Only affects users who copy-pasted from the example helpers (`example.MustSigner`, `example.MustRESTPublic`, etc.). The same setup is now inlined in each example's `main.go` using only stdlib + the public SDK. Migration is mechanical: open any current example for the equivalent setup pattern.

---

## v0.11.0

**[#106](https://github.com/amiwrpremium/go-derive/pull/106)** — typed inputs replace `map[string]any`.

Pre-v0.11, ~46 methods on `*methods.API` took `params map[string]any`. They now take typed `*Input` / `*Query` structs from `pkg/types`. Pagination is a separate `types.PageRequest` argument on methods that support it.

```go
// before
c.SendQuote(ctx, map[string]any{
    "rfq_id":               "Q-1",
    "subaccount_id":        subID,
    "max_total_cost":       cost,
    "min_total_cost":       minCost,
    // ... fields the caller had to remember from docs
})

// after
c.SendQuote(ctx, types.SendQuoteInput{
    RFQID:        "Q-1",
    SubaccountID: subID,
    MaxTotalCost: cost,
    MinTotalCost: minCost,
})
```

One behaviour change: `OrderQuotePublic` (now `GetPublicOrderQuote`) now requires a configured signer.

---

## v0.10.0

**[#104](https://github.com/amiwrpremium/go-derive/pull/104)** — pure-data records return by value.

40 method signatures change from `(*T, error)` to `(T, error)`. Field-access auto-deref means `x.Field` works on both, so most call sites don't change. Only callers doing nil checks need updating:

```go
// before
order, err := c.GetOrder(ctx, "O-1")
if err != nil || order == nil { return }
use(order.OrderID)

// after
order, err := c.GetOrder(ctx, "O-1")
if err != nil { return }
use(order.OrderID)  // err == nil means order is populated
```

Resources (`*rest.Client`, `*ws.Client`, `*ws.Subscription`, `*auth.LocalSigner`) still return pointers.

---

## v0.9.0

**[#100](https://github.com/amiwrpremium/go-derive/pull/100)** — WS subscription lifetime tied to ctx.

```go
// before — cancelling ctx aborted the subscribe RPC only;
// the subscription kept running until sub.Close.
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
sub, _ := c.SubscribeOrderBook(ctx, "BTC-PERP", "", 10)
defer sub.Close()
cancel() // sub keeps delivering

// after — cancelling ctx tears the subscription down.
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
sub, _ := c.SubscribeOrderBook(ctx, "BTC-PERP", "", 10)
defer sub.Close() // still safe, but ctx cancel also works
cancel() // sub stops, sub.Updates closes
```

To keep a subscription alive after a caller-supplied ctx expires, pass `context.Background()` and drive teardown with `sub.Close`.

---

## v0.8.x

Several breaking changes shipped through the 0.8 line.

### [#76](https://github.com/amiwrpremium/go-derive/pull/76) — input DTOs moved to `pkg/types`

```go
// before                              // after
methods.PlaceOrderInput{Asset: addr}   types.PlaceOrderInput{Asset: types.Address(addr)}
methods.GetOrdersFilter{...}           types.GetOrdersFilter{...}    // (later renamed in v0.19.0)
methods.OrderHistoryQuery{...}         types.OrderHistoryQuery{...}
methods.MMPConfig{...}                 types.MMPConfig{...}
```

`PlaceOrderInput.Asset` is now `types.Address` (the SDK wrapper) — no need to import `go-ethereum` just to construct one.

### [#78](https://github.com/amiwrpremium/go-derive/pull/78) — typed `Subscribe*` methods + wire-name drift fix

Sixteen typed convenience methods on `*ws.Client` (`SubscribeOrderBook`, `SubscribeTicker`, etc.). Generic `ws.Subscribe[T]` remains supported.

**Wire-name drift fix** — only affects code that built channel strings by hand (the typed methods always emitted the right names):

| descriptor | was | is |
|---|---|---|
| `MarginWatch` | `margin_watch` | `margin.watch` |
| `Balances` | `subaccount.{id}.balances` | `{id}.balances` |
| `Orders` | `subaccount.{id}.orders` | `{id}.orders` |
| `Quotes` | `subaccount.{id}.quotes` | `{id}.quotes` |
| `Trades` | `subaccount.{id}.trades` | `{id}.trades` |
| `TradesByTxStatus` | `subaccount.{id}.trades.{s}` | `{id}.trades.{s}` |
| `RFQs` | `wallet.{w}.rfqs` | `{w}.rfqs` |

### [#80](https://github.com/amiwrpremium/go-derive/pull/80) — `pkg/channels` removed

Channel names are plain strings; decoders are plain functions. The dedicated `pkg/channels` package is deleted. Use the typed `c.Subscribe*` methods from #78 or the generic `ws.Subscribe[T](ctx, c, channelName, decoder)`.

### [#64](https://github.com/amiwrpremium/go-derive/pull/64) — `GetOrderHistory` → `GetOrders` and back

The Go method named `GetOrderHistory` was actually calling `/private/get_orders` (status/instrument/label filtering). It is renamed to `GetOrders`, with a typed `GetOrdersFilter` (later renamed `OrdersQuery` in v0.19.0). A new `GetOrderHistory` is introduced for the time-window-based endpoint `/private/get_order_history`.

### [#63](https://github.com/amiwrpremium/go-derive/pull/63) — string fields tightened to enums

```go
// before                                  // after
auction.AuctionType == "solvent"           auction.AuctionType == enums.AuctionTypeSolvent
tx.Status == "settled"                     tx.Status == enums.TxStatusSettled
quote.InvalidReason == "..."               quote.InvalidReason == enums.RFQInvalidReason...
quote.EstimatedOrderStatus == "filled"     quote.EstimatedOrderStatus == enums.OrderStatusFilled
```

Callers comparing against the documented constants compile unchanged; callers comparing against bare string literals must switch to the constants or convert via `string(value)`.

### [#59](https://github.com/amiwrpremium/go-derive/pull/59) — `json.RawMessage` returns typed

33 methods that previously returned `json.RawMessage` now return typed DTOs from `pkg/types`. Examples:

```go
// before                                            // after
raw, err := c.GetMargin(ctx)                         marg, err := c.GetMargin(ctx)              // *types.MarginResult (later: types.MarginResult)
raw, err := c.GetAccount(ctx)                        acct, err := c.GetAccount(ctx)             // *types.AccountResult
raw, err := c.Replace(ctx, params)                   res,  err := c.Replace(ctx, params)        // *types.ReplaceResult
```

---

## v0.2.0

**[#10](https://github.com/amiwrpremium/go-derive/pull/10)** — `Trade.Realized()` removed.

```go
// before                  // after
trade.Realized()           trade.RealizedPnL
```

A field, not a method. Same value.
