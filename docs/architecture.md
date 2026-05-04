# Architecture

## The thesis

Derive's API is **JSON-RPC 2.0 over both HTTP and WebSocket** — same method
names, same params, same error envelope. Many SDKs duplicate the method
surface (`http.GetInstruments`, `ws.GetInstruments`); we don't.

Instead the SDK is layered like this:

```text
pkg/derive            top-level facade (REST + WS)
   │
   ├──► pkg/rest    ────► transport.HTTPTransport ─┐
   │     │                                          │
   │     └─ embeds ─► internal/methods.API          │
   │                                                ├─► internal/jsonrpc
   ├──► pkg/ws     ────► transport.WSTransport ────┤      (envelope + IDs)
   │     │                                          │
   │     └─ embeds ─► internal/methods.API          │
   │                                                │
   └──► pkg/channels  (typed sub descriptors)       │
                                                    │
        pkg/auth   pkg/types  pkg/enums  pkg/errors │
        ─────────────────────────────────────────── │
                                                    │
        internal/codec  internal/netconf            │
        internal/retry  internal/testutil           ┘
```

`internal/methods.API` defines each RPC method exactly once, parameterised
by `transport.Transport`. Both `pkg/rest.Client` and `pkg/ws.Client` embed
`*methods.API`, so calling `c.GetInstruments(ctx, ...)` works on either
client without code duplication.

## Why pkg/ + internal/

`pkg/` is the public API — anything users may import. Inside `internal/`
lives the plumbing (transport pumps, codec helpers, JSON-RPC framing).
Go's `internal/` rule prevents downstream users from importing it, which
lets us refactor freely.

## WebSocket lifecycle

`pkg/ws.Client` runs three goroutines under one parent context:

| Goroutine | Job |
|---|---|
| `readPump`  | Reads frames; routes to RPC futures (by id) or subscriptions (by channel) |
| `writePump` | Serialises outgoing frames (single-writer requirement of websocket libs) |
| `pingPump`  | Periodic ping; cancels on missed pong |

When the connection drops, `failConnection()` notifies all in-flight RPCs
with a `*ConnectionError` and the optional `reconnectLoop()` redials with
exponential backoff (`internal/retry.Backoff`). On a successful reconnect:

1. Re-run `OnReconnect` (the public client wires this to `Login`).
2. Re-issue `subscribe` with every channel that was alive.

User-facing subscription channels stay open across the gap. See
[reconnection.md](./reconnection.md) for the full picture.

## Auth

Two distinct signing flows, both done through the same `Signer` interface:

1. **REST/WS auth headers** — EIP-191 `personal_sign(timestamp_ms_string)`,
   sent as `X-LyraWallet`/`X-LyraTimestamp`/`X-LyraSignature`.
2. **Per-action signing** — EIP-712 `hashTypedData("Action(...)" + domain)`
   where `Action.Data = keccak(ABI(TradeModuleData))`. Signature, signer,
   nonce and expiry are embedded in the JSON-RPC params.

Two `Signer` implementations:

- `LocalSigner` — owner key in process. Owner == signer.
- `SessionKeySigner` — session key signs, but reports a separate owner
  address for the `Action.Owner` field. **Recommended for production.**

See [auth.md](./auth.md).

## Numbers

All prices, sizes, and fees use `pkg/types.Decimal` (a thin wrapper over
`shopspring/decimal`). On the wire they are JSON strings — Derive's
convention so that `1e-18` precision survives round-trips. For ABI-encoded
action signing, `internal/codec.DecimalToU256` and `DecimalToI256` scale to
1e18 fixed-point. See [numerics.md](./numerics.md).

## Errors

Three-level model, all `errors.Is` / `errors.As` compatible:

| Layer | Type | Use |
|---|---|---|
| Sentinels | `var ErrUnauthorized = errors.New(...)` | Compare with `errors.Is` for known kinds |
| API errors | `*APIError{Code,Message,Data}` | JSON-RPC error from server; `Is` maps codes back to sentinels |
| Network | `*ConnectionError`, `*TimeoutError` | Wrap transport-level failures; `Unwrap` exposes the original |

The full catalogue of 136 server-side codes lives in
[error-handling.md](./error-handling.md).

## Numeric scale of the codebase

| Surface | Status |
|---|---|
| Public packages | 10 (`pkg/derive`, `rest`, `ws`, `auth`, `types`, `enums`, `errors`, `channels`, `channels/public`, `channels/private`, `contracts`) |
| Internal packages | 7 (`jsonrpc`, `transport`, `methods`, `netconf`, `codec`, `retry`, `testutil`) |
| Source files (`pkg/` + `internal/`) | 101 |
| Unit-test files | 109 |
| `Example_*` functions in package tests | 0 — examples live under `examples/` |
| Fuzz test files | 8 (`Fuzz*` functions: 10) |
| Runnable example programs | 91, one per directory under `examples/` |
| Integration test files | 7 under `test/`, gated by `-tags=integration` |
| CI workflows | 24 (ci, lint, extra-lint, codeql, gosec, semgrep, codacy, scorecard, osv-scanner, trivy, gitleaks, trufflehog, dependency-review, license-check, pin-check, release, release-please, verify-release, integration, pr-title, labeler, auto-merge, auto-assign, stale) |

## Methods not exposed

The SDK covers every JSON-RPC method that Derive exposes for trading,
querying account state, and streaming market data. The endpoints listed
below are **deliberately not wrapped** — they all require on-chain
signing of an actual transaction (token transfer, contract registration,
admin call), which is out of scope for the JSON-RPC surface. Callers
needing them should compose the on-chain call themselves and submit it
via Derive's contracts directly.

| Method | Why it's out of scope here |
|---|---|
| `private/deposit`, `private/withdraw` | ERC-20 transfer to/from the Derive bridge — needs go-ethereum bindings + signed tx |
| `private/transfer_erc20`, `private/transfer_position` | inter-subaccount on-chain transfer |
| `private/create_subaccount` | on-chain subaccount registration |
| `private/session_keys`, `private/change_session_key_label` | session-key registry plumbing |
| `public/register_session_key`, `public/deregister_session_key`, `public/build_register_session_key_tx` | session-key registration helpers (admin-side) |
| `public/create_account`, `public/create_subaccount_debug` | account bootstrap; admin/debug |
| `public/change_compliance_status`, `public/set_feed_data`, `public/margin_watch` | admin / oracle-write paths |
