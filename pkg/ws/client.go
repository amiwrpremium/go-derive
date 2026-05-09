// Package ws is the WebSocket-backed client for Derive's JSON-RPC API.
//
// # What it covers
//
// Derive's WebSocket transport carries two distinct workloads:
//
//   - request/response RPCs (lower latency than HTTP because of connection
//     reuse and no per-call TLS handshake)
//   - pub/sub channel notifications (the only way to stream live data)
//
// [Client] handles both. It runs three goroutines under one parent context:
// a read pump that demultiplexes responses from notifications, a write pump
// that serialises outgoing frames, and a ping pump that keeps the connection
// alive. When [WithReconnect] is enabled, a reconnect goroutine re-dials
// with exponential backoff and re-issues subscribe + login on success.
//
// # Lifecycle
//
//	c, _ := ws.New(ws.WithMainnet(), ws.WithSigner(s), ws.WithSubaccount(123))
//	defer c.Close()
//	if err := c.Connect(ctx); err != nil { ... }
//	if err := c.Login(ctx); err != nil { ... }
//
//	sub, err := ws.Subscribe[types.OrderBook](ctx, c, public.OrderBook{Instrument: "BTC-PERP"})
//	defer sub.Close()
//	for ob := range sub.Updates() { ... }
//
// # Concurrency
//
// [Client] is safe for concurrent use after Connect. Many goroutines may
// call methods or hold subscriptions on the same client simultaneously.
package ws

import (
	"context"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/transport"
)

// Client is a JSON-RPC plus subscription client over a single WebSocket.
//
// Construct with [New], then [Client.Connect] and (for private endpoints)
// [Client.Login]. The zero value is not usable.
type Client struct {
	*derive.API
	wt     *transport.WSTransport
	signer derive.Signer
	cfg    derive.NetworkConfig
}

// New constructs a [Client] without dialing — Connect opens the socket.
//
// One of [WithMainnet], [WithTestnet], or [WithCustomNetwork] must be
// supplied. All other options have sensible defaults: 20-second ping
// interval, 10 TPS rate limit with 5x burst, auto-reconnect enabled,
// 5-minute signature expiry.
func New(opts ...Option) (*Client, error) {
	c := &config{
		userAgent:    derive.UserAgent(),
		expiry:       300,
		tps:          10,
		burst:        5,
		pingInterval: 20 * time.Second,
		reconnect:    true,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.network.Network == derive.NetworkUnknown {
		return nil, derive.ErrInvalidConfig
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

	api := &derive.API{
		T:               wt,
		Signer:          c.signer,
		Domain:          c.network.EIP712Domain(),
		Subaccount:      c.subaccount,
		Nonces:          derive.NewNonceGen(),
		SignatureExpiry: c.expiry,
	}
	api.SetTradeModule(common.HexToAddress(c.network.Contracts.TradeModule))

	cli := &Client{API: api, wt: wt, signer: c.signer, cfg: c.network}

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
// channels and private RPCs. Call once after [Client.Connect]; the
// reconnect goroutine re-issues it on every redial when [WithReconnect]
// is enabled.
//
// Returns [github.com/amiwrpremium/go-derive/pkg/errors.ErrUnauthorized]
// if no signer was configured. Server-side rejections come through as
// [github.com/amiwrpremium/go-derive/pkg/errors.APIError].
//
// Login is WS-only — Derive's REST surface authenticates per-request
// via signed headers, not via this RPC.
func (c *Client) Login(ctx context.Context) error {
	if c.signer == nil {
		return derive.ErrUnauthorized
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
		// Transport returns *transport.JSONRPCError; lift to the
		// public *derive.APIError at the boundary so callers can
		// match `errors.As(err, &derive.APIError{...})`.
		if rpcErr, ok := err.(*transport.JSONRPCError); ok {
			return &derive.APIError{Code: rpcErr.Code, Message: rpcErr.Message, Data: rpcErr.Data}
		}
		return err
	}
	return nil
}

// installReconnectLogin tells the underlying transport to re-issue the
// `public/login` RPC each time the socket is re-established.
func (c *Client) installReconnectLogin() {
	c.wt.SetOnReconnect(func(ctx context.Context, _ *transport.WSTransport) error {
		return c.Login(ctx)
	})
}

// Connect dials the WebSocket and starts the read/write/ping pumps.
//
// It returns once the handshake completes or fails. Errors are wrapped in
// [github.com/amiwrpremium/go-derive/pkg/errors.ConnectionError].
func (c *Client) Connect(ctx context.Context) error { return c.wt.Connect(ctx) }

// Close terminates the WebSocket, unblocks any in-flight calls with
// [github.com/amiwrpremium/go-derive/pkg/errors.ErrSubscriptionClosed], and
// stops the pump goroutines. Close is idempotent. The Client is unusable
// after Close.
func (c *Client) Close() error { return c.wt.Close() }

// IsConnected reports whether the underlying socket is currently open.
func (c *Client) IsConnected() bool { return c.wt.IsConnected() }

// Network returns the active network configuration.
func (c *Client) Network() derive.NetworkConfig { return c.cfg }

// transport exposes the WS transport to the subscription helpers in this
// package. Lower-cased so it stays internal to pkg/ws.
func (c *Client) transport() *transport.WSTransport { return c.wt }
