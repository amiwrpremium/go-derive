# Error handling

Three error layers, all `errors.Is` / `errors.As` friendly.

## Layer 1: sentinels

Defined in `errors.go/errors.go`. Use `errors.Is` to compare:

```go
if errors.Is(err, derive.ErrRateLimited) {
    backoff()
    return
}
```

The 21 sentinels:

- `ErrNotConnected`, `ErrAlreadyConnected`
- `ErrUnauthorized`, `ErrInvalidSignature`, `ErrSessionKeyExpired`,
  `ErrSessionKeyNotFound`
- `ErrRateLimited`
- `ErrInsufficientFunds`, `ErrOrderNotFound`, `ErrAlreadyCancelled`,
  `ErrAlreadyFilled`, `ErrAlreadyExpired`
- `ErrInstrumentNotFound`, `ErrSubaccountNotFound`, `ErrAccountNotFound`
- `ErrChainIDMismatch`, `ErrMMPFrozen`, `ErrRestrictedRegion`
- `ErrSubscriptionClosed`, `ErrSubaccountRequired`, `ErrInvalidConfig`

## Layer 2: `*APIError`

Every server-side JSON-RPC error becomes a `*APIError`:

```go
type APIError struct {
    Code    int             // see codes.go
    Message string          // server-supplied human message
    Data    json.RawMessage // optional structured detail
}
```

Use `errors.As` to inspect:

```go
var apiErr *derive.APIError
if errors.As(err, &apiErr) {
    log.Printf("derive code %d: %s", apiErr.Code, apiErr.CanonicalMessage())
}
```

`APIError.Is` maps codes back to sentinels — so `errors.Is(err,
ErrRateLimited)` works regardless of whether the underlying error is the
raw `*APIError` or a sentinel.

## Layer 3: transport errors

`*ConnectionError` and `*TimeoutError` wrap network-level failures with
`Unwrap` so the original `net.Error` (or context error) is reachable.

```go
var connErr *derive.ConnectionError
if errors.As(err, &connErr) {
    log.Printf("transport: %s: %v", connErr.Op, connErr.Err)
}
```

## The 136 server codes

The full catalogue is in `errors.go/codes.go`, grouped by topic:

| Range | Topic |
|---|---|
| `-32700`…`-32603` | JSON-RPC 2.0 standard |
| `-32000`, `-32100` | rate limiting |
| `9000`–`9001` | engine timeouts |
| `10000`–`10015` | account / wallet |
| `11000`–`11021` | order placement / lifecycle |
| `11050`–`11055` | trigger orders |
| `11100`–`11107` | RFQ / quotes |
| `11200`–`11203` | auctions |
| `12000`–`12003` | assets / instruments |
| `13000` | subscriptions |
| `14000`–`14034` | account / auth |
| `16000`, `16001`, `16100` | compliance |
| `18000`–`18007` | vault / block |
| `19000` | maker programs |

## Sentinel → code map

```go
var apiErr *derive.APIError
errors.As(err, &apiErr)
// apiErr.Code lookup ↓
```

| Sentinel | Codes that match |
|---|---|
| `ErrRateLimited` | `-32000`, `-32100` |
| `ErrUnauthorized` | `14014`, `14020`, `14021`, `14023`, `14025`, `14026`, `14027`, `14029`, `14030`, `14031`, `14032`, `14033`, `16100` |
| `ErrInvalidSignature` | `14014` |
| `ErrSessionKeyExpired` | `14030` |
| `ErrSessionKeyNotFound` | `14026` |
| `ErrInsufficientFunds` | `11000`, `10011`, `10012` |
| `ErrOrderNotFound` | `11006` |
| `ErrAlreadyCancelled/Filled/Expired` | `11003` / `11004` / `11005` |
| `ErrInstrumentNotFound` | `12001`, `12000` |
| `ErrSubaccountNotFound` | `14001` |
| `ErrAccountNotFound` | `14000` |
| `ErrChainIDMismatch` | `14024` |
| `ErrMMPFrozen` | `11015` |
| `ErrRestrictedRegion` | `16000`, `16001` |

Codes outside this map don't match any sentinel — drill in via
`errors.As(&apiErr)` and inspect `apiErr.Code` directly.

## Canonical messages

Each code has a human-readable description in `errors.go/messages.go`.
When the server returns a sparse `Message`, `APIError.Error()` falls back
to the canonical text:

```go
e := &derive.APIError{Code: derive.CodeMMPFrozen}
fmt.Println(e.Error())
// derive: api error 11015: market-maker protection has frozen this currency
```

You can also look it up explicitly:

```go
desc := derive.Description(derive.CodeSessionKeyExpired)
// "session key has expired"
```

## A retry pattern

```go
const maxRetries = 5
backoff := 200 * time.Millisecond

for i := 0; i < maxRetries; i++ {
    err := c.PlaceOrder(ctx, in)
    switch {
    case err == nil:
        return nil
    case errors.Is(err, derive.ErrRateLimited):
        time.Sleep(backoff)
        backoff *= 2
    case errors.Is(err, derive.ErrSessionKeyExpired):
        // re-register session key (operator action) and abort
        return err
    case errors.Is(err, derive.ErrInsufficientFunds):
        return err // not retriable
    default:
        // unknown — bubble up
        return err
    }
}
```

The signing path's own errors (key parsing, hashing) come back as
`*SigningError`; the order-expiry guard surfaces as
`*ExpiredSignatureError` — both with `Unwrap` for the cause.
