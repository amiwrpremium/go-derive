# WebSocket reconnection

The WebSocket transport (`internal/transport/ws.go`) is the most
intricate piece in the SDK. This page is a guided tour for anyone
debugging a connection issue or extending the lifecycle.

## Three pumps under one connection

When `Connect` succeeds, three goroutines start under one parent context:

| Pump | Reads from | Writes to | Exit signal |
|---|---|---|---|
| `readPump` | `conn.Read` | dispatches to RPC futures or subscription channels | conn closes |
| `writePump` | `writeQ` (chan []byte) | `conn.Write` | `stopCh` closed |
| `pingPump` | timer (default 20s) | `conn.Ping` | `stopCh` closed |

Each pump receives its `conn`, `wq`, and `stop` as function arguments
from `dial()` — *not* via `t.conn` field reads. That avoids the nil-conn
race we hit and fixed during the unit-test work.

## ID correlation

Every outgoing RPC gets a fresh ID from `internal/jsonrpc.IDGen`. The
`Call` method puts a `pendingCall{id → chan err}` entry in the transport's
map and writes the request to `writeQ`. `readPump` looks up incoming
responses by id, removes the map entry, and signals via the `pendingCall`
channel.

Subscription notifications have no id (`IsNotification` peeks at the
JSON to distinguish them). They're routed by `params.channel` to the
matching `*wsSub`'s update channel.

## Failure path

When `conn.Read` or `conn.Write` returns a non-nil error,
`failConnection(err)` runs:

1. Snapshot and clear the pending map.
2. Set `t.conn = nil`, cancel the root context, close `stopCh`.
3. Notify every pending RPC with `*ConnectionError{Op, Err}`.
4. Best-effort `conn.Close(StatusNormalClosure, "fail")`.

If `Reconnect` is enabled, the `reconnectLoop()` goroutine picks up from
here.

## Reconnect loop

```go
func reconnectLoop():
    bo := retry.NewBackoff()  // exp backoff, max 30s
    for {
        if t.conn == nil:
            ctx := timeout(30s)
            if err := dial(ctx); err != nil:
                sleep(bo.Next())
                continue
            bo.Reset()
            if cfg.OnReconnect != nil:
                _ = cfg.OnReconnect(timeout(30s), t)  // re-login
            resubscribe()  // re-issue subscribe RPCs for active subs
        sleep(1s)
    }
```

`pkg/ws.Client.installReconnectLogin` wires `OnReconnect` to `c.Login`
when a signer is configured, so private subscriptions resume cleanly.

## What stays open

User-facing `*Subscription[T]` channels are *not* closed during a
reconnect — that's the whole point. The transport's underlying
`transport.Subscription` keeps the same `updates` channel; the new
post-reconnect WS frames flow into it without the consumer noticing.

If a subscription is explicitly closed by the user (`sub.Close()`), it's
removed from the active set and won't be re-subscribed.

## Pings and pongs

`pingPump` sends an application-level ping every 20s by default. If the
server doesn't respond within 5s, `conn.Ping` returns an error,
`failConnection` runs, and reconnect kicks in. The interval is tunable:

```go
ws.New(ws.WithTestnet(), ws.WithPingInterval(10*time.Second))
```

## Disabling reconnect

```go
ws.New(ws.WithTestnet(), ws.WithReconnect(false))
```

With reconnect off, a connection drop returns `*ConnectionError` to every
pending RPC and closes every subscription's update channel. The caller
is responsible for handling the close and dialing again.

## Edge cases handled in tests

- **Server hangs up mid-call**: pending RPCs receive
  `*ConnectionError{Op: "ws read", Err: ...}`. Test:
  `TestWSTransport_FailureClosesPending`.
- **Subscribe to the same channel twice**: the second call returns the
  same handle, no duplicate RPC. Test:
  `TestWSTransport_SubscribeIdempotent`.
- **Server-side subscribe error**: the local subscription record is
  freed before the error propagates. Test:
  `TestWSTransport_SubscribeServerError`.
- **Race between `Close` and pingPump first tick**: each pump owns its
  conn pointer so a half-cleared transport state can't yield a nil
  dereference. Caught by `-race` mode.

## Implementation files

| File | What's in it |
|---|---|
| `internal/transport/ws.go` | the transport itself |
| `internal/transport/ws_test.go` | RPC + subscription tests against a mock WS server |
| `internal/transport/ws_reconnect_test.go` | reconnect/idempotency tests |
| `internal/testutil/mockws.go` | the in-process mock server |
| `internal/retry/backoff.go` | exponential backoff with jitter |
