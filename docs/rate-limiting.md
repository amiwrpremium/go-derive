# Rate limiting

Every transport (REST and WS) holds an independent token-bucket limiter.

## Defaults

| Knob | Default | Why |
|---|---|---|
| sustained TPS | 10 | Derive's documented per-IP limit |
| burst | 5× | Derive's documented per-IP burst |
| capacity | tps × burst = 50 | tokens at full bucket |
| refill rate | 10 tokens/sec | sustained = TPS |

A nil `*RateLimiter` is treated as "limiting disabled" — every operation
is a no-op, never panics. `NewRateLimiter(0, *)` returns nil; that's the
shortcut for opting out.

## Tuning

```go
c, _ := derive.NewRestClient(
    derive.WithTestnet(),
    derive.WithRateLimit(50, 2), // 50 TPS, 100-token bucket
)
```

The same option exists for `derive.WsClient`:

```go
c, _ := derive.NewWsClient(derive.WithTestnet(), derive.WithRateLimit(50, 2))
```

## Behaviour under saturation

When the bucket is empty, `RateLimiter.Wait(ctx)` blocks until a token
refills *or* `ctx` cancels — whichever comes first. The transport's
`Call(ctx, ...)` calls `Wait` first, so a saturated rate limit shows up
as a context-cancellation error rather than a server-side
`ErrRateLimited`.

To distinguish in your retry logic:

```go
switch {
case errors.Is(err, context.DeadlineExceeded):
    // local rate limit hit, ctx expired before token refilled
case errors.Is(err, derive.ErrRateLimited):
    // server told us to back off (we lost a race with another caller)
}
```

## REST + WS share nothing

REST and WS limit independently. If you saturate WS, REST still has its
own headroom. If you need a single combined budget, build a top-level
limiter and call its `Wait` before every dispatch.

## Production recommendations

- **Don't disable** the limiter unless you've already proxied behind your
  own gateway. Hitting Derive's IP limit returns code `-32000` and
  triggers a temporary block.
- **Match `tps` to your account's actual quota**; bursts above sustained
  trigger soft rate-limit responses you can recover from, but sustained
  over-limit triggers harder blocks.
- **Different processes need different limiters.** Two SDK clients in
  different processes don't coordinate. If you're running many trading
  bots from one IP, give each its own slice of the budget.

## Internals

The implementation is the standard token-bucket: `Wait` reads, refills
based on elapsed time, and either decrements + returns immediately or
sleeps for the time-until-next-token.

```go
// pseudocode of internal/transport.RateLimiter.Wait
for {
    if available_tokens >= 1:
        take_token()
        return nil
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(time_until_next_token):
        // continue, recompute
    }
}
```

The full implementation lives in `internal/transport/ratelimit.go` and is
covered by `internal/transport/ratelimit_test.go` including blocking,
ctx-cancel, and burst behaviour.
