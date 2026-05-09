// WebSocket-backed JSON-RPC client lives in this file: [WsClient], the
// generic [Subscribe]/[SubscribeFunc] helpers, and the WS-specific
// [WithPingInterval] / [WithReconnect] knobs. The shared [Option]
// type and the network/signer/subaccount/rate-limit/user-agent
// builders live in rest.go and are reused unchanged.
//
// # What it covers
//
// Derive's WebSocket transport carries two distinct workloads:
//
//   - request/response RPCs (lower latency than HTTP because of
//     connection reuse and no per-call TLS handshake)
//   - pub/sub channel notifications (the only way to stream live data)
//
// [WsClient] handles both. It runs three goroutines under one parent
// context: a read pump that demultiplexes responses from notifications,
// a write pump that serialises outgoing frames, and a ping pump that
// keeps the connection alive. When [WithReconnect] is enabled, a
// reconnect goroutine re-dials with exponential backoff and re-issues
// subscribe + login on success.
//
// # Lifecycle
//
//	c, _ := derive.NewWsClient(derive.WithMainnet(), derive.WithSigner(s), derive.WithSubaccount(123))
//	defer c.Close()
//	if err := c.Connect(ctx); err != nil { ... }
//	if err := c.Login(ctx); err != nil { ... }
//
//	sub, err := derive.Subscribe[derive.OrderBook](ctx, c, derive.PublicOrderBook{Instrument: "BTC-PERP"})
//	defer sub.Close()
//	for ob := range sub.Updates() { ... }
//
// # Concurrency
//
// [WsClient] is safe for concurrent use after Connect. Many goroutines
// may call methods or hold subscriptions on the same client
// simultaneously.

package derive

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/internal/transport"
)

// WsClient is a JSON-RPC plus subscription client over a single
// WebSocket.
//
// Construct with [NewWsClient], then [WsClient.Connect] and (for
// private endpoints) [WsClient.Login]. The zero value is not usable.
type WsClient struct {
	*API
	wt     *transport.WSTransport
	signer Signer
	cfg    NetworkConfig
}

// NewWsClient constructs a [WsClient] without dialing — Connect opens
// the socket.
//
// One of [WithMainnet], [WithTestnet], or [WithCustomNetwork] must be
// supplied. All other options have sensible defaults: 20-second ping
// interval, 10 TPS rate limit with 5x burst, auto-reconnect enabled,
// 5-minute signature expiry.
func NewWsClient(opts ...Option) (*WsClient, error) {
	c := &config{
		userAgent:    UserAgent(),
		expiry:       300,
		tps:          10,
		burst:        5,
		pingInterval: 20 * time.Second,
		reconnect:    true,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.network.Network == NetworkUnknown {
		return nil, ErrInvalidConfig
	}

	wt, err := transport.NewWS(transport.WSConfig{
		URL:          c.network.WSURL,
		Limiter:      transport.NewRateLimiter(c.tps, c.burst),
		UserAgent:    c.userAgent,
		PingInterval: c.pingInterval,
		Reconnect:    c.reconnect,
	})
	if err != nil {
		return nil, err
	}

	api := &API{
		T:               wt,
		Signer:          c.signer,
		Domain:          c.network.EIP712Domain(),
		Subaccount:      c.subaccount,
		Nonces:          NewNonceGen(),
		SignatureExpiry: c.expiry,
	}
	api.SetTradeModule(common.HexToAddress(c.network.Contracts.TradeModule))

	cli := &WsClient{API: api, wt: wt, signer: c.signer, cfg: c.network}

	// Re-login automatically after each reconnect so private subscriptions
	// resume working. Skipped when no signer is configured.
	if c.reconnect && c.signer != nil {
		cli.installReconnectLogin()
	}
	return cli, nil
}

// Login authenticates the WebSocket session with the configured signer.
//
// Sends an EIP-191 personal-sign over the millisecond timestamp and
// dispatches the `public/login` RPC, which gates access to the private
// channels and private RPCs. Call once after [WsClient.Connect]; the
// reconnect goroutine re-issues it on every redial when [WithReconnect]
// is enabled.
//
// Returns [ErrUnauthorized] if no signer was configured. Server-side
// rejections come through as [*APIError].
//
// Login is WS-only — Derive's REST surface authenticates per-request
// via signed headers, not via this RPC.
func (c *WsClient) Login(ctx context.Context) error {
	if c.signer == nil {
		return ErrUnauthorized
	}
	now := time.Now()
	sig, err := c.signer.SignAuthHeader(ctx, now)
	if err != nil {
		return err
	}
	params := map[string]any{
		"wallet":    c.signer.Owner().Hex(),
		"timestamp": strconv.FormatInt(now.UnixMilli(), 10),
		"signature": sig.Hex(),
	}
	if err := c.wt.Call(ctx, "public/login", params, nil); err != nil {
		// Transport returns *transport.JSONRPCError; lift to the public
		// *APIError at the boundary so callers can match
		// `errors.As(err, &APIError{...})`.
		if rpcErr, ok := err.(*transport.JSONRPCError); ok {
			return &APIError{Code: rpcErr.Code, Message: rpcErr.Message, Data: rpcErr.Data}
		}
		return err
	}
	return nil
}

// installReconnectLogin tells the underlying transport to re-issue the
// `public/login` RPC each time the socket is re-established.
func (c *WsClient) installReconnectLogin() {
	c.wt.SetOnReconnect(func(ctx context.Context, _ *transport.WSTransport) error {
		return c.Login(ctx)
	})
}

// Connect dials the WebSocket and starts the read/write/ping pumps.
//
// It returns once the handshake completes or fails. Errors are wrapped
// in [*ConnectionError].
func (c *WsClient) Connect(ctx context.Context) error { return c.wt.Connect(ctx) }

// Close terminates the WebSocket, unblocks any in-flight calls with
// [ErrSubscriptionClosed], and stops the pump goroutines. Close is
// idempotent. The WsClient is unusable after Close.
func (c *WsClient) Close() error { return c.wt.Close() }

// IsConnected reports whether the underlying socket is currently open.
func (c *WsClient) IsConnected() bool { return c.wt.IsConnected() }

// Network returns the active network configuration.
func (c *WsClient) Network() NetworkConfig { return c.cfg }

// rawTransport exposes the WS transport to the subscription helpers.
// Lower-cased so it stays internal to the package.
func (c *WsClient) rawTransport() *transport.WSTransport { return c.wt }

// WithPingInterval overrides the periodic application-level ping
// cadence. The default is 20 seconds. Setting it too high risks the
// connection being pruned by intermediate proxies.
//
// WS-only — ignored by [RestClient].
func WithPingInterval(d time.Duration) Option { return func(c *config) { c.pingInterval = d } }

// WithReconnect controls whether the client auto-reconnects after the
// underlying connection drops. When enabled (the default), the client
// re-dials with exponential backoff, re-runs login (if a signer is
// configured), and re-issues every active subscription so user-facing
// channels stay open across the gap.
//
// WS-only — ignored by [RestClient].
func WithReconnect(enabled bool) Option { return func(c *config) { c.reconnect = enabled } }

// Subscribe registers a typed subscription on a [WsClient] and returns
// a [Subscription] whose Updates channel yields values of type T.
//
// T must match the type the channel descriptor's Decode method
// returns; a mismatch is dropped silently rather than crashing the
// read pump (the underlying decoder error is surfaced if a debugger
// is attached). Pass the right T for the descriptor — e.g. [OrderBook]
// for [PublicOrderBook], [][Trade] for [PublicTrades].
//
// Generics let callers avoid type assertions at the use site:
//
//	sub, _ := derive.Subscribe[derive.OrderBook](ctx, c,
//	    derive.PublicOrderBook{Instrument: "BTC-PERP"})
//	defer sub.Close()
//	for ob := range sub.Updates() {
//	    fmt.Println(ob.Bids[0])
//	}
//
// The returned subscription buffers up to 256 events in memory; if
// the caller is slow, newer events are dropped (best-effort fan-out,
// not a reliable queue). Use [SubscribeFunc] when you want to be sure
// every event is processed.
func Subscribe[T any](ctx context.Context, c *WsClient, ch Channel) (*Subscription[T], error) {
	dec := func(raw json.RawMessage) (any, error) {
		v, err := ch.Decode(raw)
		if err != nil {
			return nil, err
		}
		typed, ok := v.(T)
		if !ok {
			return nil, fmt.Errorf("derive: channel %q: decoded type %T does not match expected %T", ch.Name(), v, *new(T))
		}
		return typed, nil
	}
	sub, err := c.rawTransport().Subscribe(ctx, ch.Name(), dec)
	if err != nil {
		return nil, err
	}
	out := &Subscription[T]{
		raw:     sub,
		typed:   make(chan T, 256),
		channel: ch.Name(),
	}
	go out.pump()
	return out, nil
}

// SubscribeFunc is a convenience over [Subscribe] that drives a
// per-event callback synchronously. It returns when ctx is cancelled
// (returning ctx.Err()) or the subscription closes (returning the
// underlying terminal error, which may be nil for a clean close).
//
// Use SubscribeFunc when callback-driven code reads more naturally
// than a channel-receive loop, or when you want to guarantee every
// event is processed (the callback runs synchronously, so back-pressure
// on the caller is back-pressure on the subscription).
func SubscribeFunc[T any](ctx context.Context, c *WsClient, ch Channel, fn func(T)) error {
	sub, err := Subscribe[T](ctx, c, ch)
	if err != nil {
		return err
	}
	defer func() { _ = sub.Close() }()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-sub.Updates():
			if !ok {
				return sub.Err()
			}
			fn(ev)
		}
	}
}

// Subscription is a typed wrapper around the underlying transport-level
// subscription. The zero value is not usable; obtain one from
// [Subscribe].
//
// Always call [Subscription.Close] to release the channel slot and
// tell the server to stop sending updates. The Close call is
// idempotent.
type Subscription[T any] struct {
	raw     transport.Subscription
	typed   chan T
	channel string
}

// Channel returns the dotted server-side channel name (e.g.
// "orderbook.BTC-PERP.1.10"). Useful for diagnostics and logs.
func (s *Subscription[T]) Channel() string { return s.channel }

// Updates returns the receive channel of typed events. The channel is
// closed when the subscription terminates; receivers should select
// against ctx.Done() to know when to bail.
func (s *Subscription[T]) Updates() <-chan T { return s.typed }

// Err returns the terminal error once [Subscription.Updates] has
// closed, or nil for a clean shutdown.
func (s *Subscription[T]) Err() error { return s.raw.Err() }

// Close ends the subscription, sends an unsubscribe RPC best-effort,
// and drains the typed channel. Idempotent.
func (s *Subscription[T]) Close() error { return s.raw.Close() }

// pump bridges the untyped transport channel to the typed user
// channel. Type-mismatched events are dropped (Subscribe returns an
// error if T can't accept the descriptor's output, but we still
// defend at runtime).
func (s *Subscription[T]) pump() {
	defer close(s.typed)
	for v := range s.raw.Updates() {
		typed, ok := v.(T)
		if !ok {
			continue
		}
		s.typed <- typed
	}
}
