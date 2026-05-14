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
	"time"

	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

// Option configures a [Client] at construction time. Compose any number of
// With* helpers and pass them to [New].
type Option func(*config)

type config struct {
	network         netconf.Config
	signer          auth.Signer
	subaccount      int64
	userAgent       string
	tps             float64
	burst           float64
	pingInterval    time.Duration
	reconnect       bool
	expiry          int64
	preloadAllInsts bool
}

// WithMainnet selects Derive's mainnet endpoints (chain id 957).
func WithMainnet() Option { return func(c *config) { c.network = netconf.Mainnet() } }

// WithTestnet selects Derive's demo (testnet) endpoints (chain id 901).
func WithTestnet() Option { return func(c *config) { c.network = netconf.Testnet() } }

// WithCustomNetwork overrides the entire network configuration. Use it for
// staging or vendored deployments.
func WithCustomNetwork(cfg netconf.Config) Option { return func(c *config) { c.network = cfg } }

// WithSigner attaches an auth [github.com/amiwrpremium/go-derive/pkg/auth.Signer]
// used for both per-action EIP-712 signing and the WS `public/login`
// timestamp signature.
//
// Without a signer, only public RPCs and channels work; private calls
// return [github.com/amiwrpremium/go-derive/pkg/errors.ErrUnauthorized].
func WithSigner(s auth.Signer) Option { return func(c *config) { c.signer = s } }

// WithSubaccount sets the subaccount id used by private methods and
// private subscription channels.
func WithSubaccount(id int64) Option { return func(c *config) { c.subaccount = id } }

// WithUserAgent overrides the User-Agent header sent on the WebSocket
// upgrade request. Default is "go-derive/<version>".
func WithUserAgent(ua string) Option { return func(c *config) { c.userAgent = ua } }

// WithRateLimit configures the per-client token-bucket rate limiter.
//
// tps is the sustained transactions-per-second; burst sets the bucket
// capacity (capacity = tps × burst). Pass tps <= 0 to disable. Defaults:
// 10 TPS, burst 5.
func WithRateLimit(tps, burst float64) Option {
	return func(c *config) {
		c.tps = tps
		c.burst = burst
	}
}

// WithPingInterval overrides the periodic application-level ping cadence.
// The default is 20 seconds. Setting it too high risks the connection being
// pruned by intermediate proxies.
func WithPingInterval(d time.Duration) Option { return func(c *config) { c.pingInterval = d } }

// WithReconnect controls whether the client auto-reconnects after the
// underlying connection drops. When enabled (the default), the client
// re-dials with exponential backoff, re-runs login (if a signer is
// configured), and re-issues every active subscription so user-facing
// channels stay open across the gap.
func WithReconnect(enabled bool) Option { return func(c *config) { c.reconnect = enabled } }

// WithSignatureExpiry sets the seconds-from-now expiry on signed actions.
// The default is 300 (5 minutes).
func WithSignatureExpiry(seconds int64) Option { return func(c *config) { c.expiry = seconds } }

// WithInstrumentPreload kicks off a background fetch of every live
// instrument (across all currencies and kinds) immediately after the
// client is constructed.
//
// The preload populates the SDK's instrument metadata cache, letting
// subsequent signed actions (PlaceOrder, SendQuote, etc.) skip the
// per-instrument public/get_instrument lookup that otherwise happens
// lazily on first use of each new instrument.
//
// The fetch runs in a goroutine using context.Background(); errors
// are swallowed. Callers who want to surface errors should invoke
// [methods.API.PreloadAllInstruments] manually instead.
func WithInstrumentPreload() Option {
	return func(c *config) { c.preloadAllInsts = true }
}
