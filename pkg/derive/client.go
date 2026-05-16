// Package derive is the top-level convenience facade for the SDK. It
// composes a [github.com/amiwrpremium/go-derive/pkg/rest.Client] and a
// [github.com/amiwrpremium/go-derive/pkg/ws.Client] sharing the same
// signer, subaccount and network configuration.
//
// # Most users want this
//
//	c, _ := derive.NewClient(
//	    derive.WithMainnet(),
//	    derive.WithSigner(signer),
//	    derive.WithSubaccount(123),
//	)
//	defer c.Close()
//
//	insts, _ := c.REST.GetInstruments(ctx, types.InstrumentsQuery{Currency: "BTC", Kind: enums.InstrumentTypePerp})
//	c.WS.Connect(ctx)
//	c.WS.Login(ctx)
//
// # Option parity with the transports
//
// Every option that pkg/rest or pkg/ws accepts is also exposed here.
// Common options ([WithSigner], [WithSubaccount], [WithUserAgent],
// [WithRateLimit], [WithSignatureExpiry], [WithInstrumentPreload]) flow
// to both transports. Transport-specific options route to whichever
// transport supports them: [WithHTTPClient] and [WithHTTPTimeout] go to
// REST only; [WithPingInterval], [WithReconnect], and [WithOnReconnect]
// go to WS only.
//
// # When to skip the facade
//
// Use pkg/rest or pkg/ws directly when you only need one transport —
// both expose the full RPC method surface independently. The facade is
// a shortcut for the common case where you want both bundled under
// one signer + subaccount.
package derive

import (
	"context"
	"net/http"
	"time"

	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/pkg/auth"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

// Client bundles a REST client and a WebSocket client sharing the same
// signer, subaccount and network configuration.
//
// Both REST and WS are exported so callers can reach the underlying client
// for transport-specific options (e.g. c.WS.Connect(ctx)). The zero value
// is not usable; obtain one from [NewClient].
type Client struct {
	// REST is the HTTP-backed JSON-RPC client.
	REST *rest.Client
	// WS is the WebSocket-backed JSON-RPC + subscription client.
	WS *ws.Client

	cfg netconf.Config
}

// Option configures a [Client] at construction time. Compose any number of
// With* helpers and pass them to [NewClient].
type Option func(*config)

// config carries the facade-side state. Most fields are nilable so the
// facade can distinguish "user explicitly set this" from "user left the
// default" — when unset, the transport's own default applies, so we
// don't have to track the transport defaults here (and risk drifting
// from them).
type config struct {
	network      netconf.Config
	signer       auth.Signer
	subaccount   int64
	connectWS    bool
	preloadInsts bool

	// Common options (forwarded to both REST and WS when set).
	userAgent string // "" means unset
	tps       *float64
	burst     *float64
	expiry    *int64

	// REST-only options.
	httpClient  *http.Client
	httpTimeout *time.Duration

	// WS-only options.
	pingInterval *time.Duration
	reconnect    *bool
	onReconnect  func(error)
}

// WithMainnet selects Derive's mainnet endpoints.
func WithMainnet() Option { return func(c *config) { c.network = netconf.Mainnet() } }

// WithTestnet selects Derive's demo (testnet) endpoints.
func WithTestnet() Option { return func(c *config) { c.network = netconf.Testnet() } }

// WithCustomNetwork overrides the entire network configuration. Use it for
// staging or vendored deployments.
func WithCustomNetwork(cfg netconf.Config) Option { return func(c *config) { c.network = cfg } }

// WithSigner attaches an auth [github.com/amiwrpremium/go-derive/pkg/auth.Signer].
// Without one, only public endpoints work.
func WithSigner(s auth.Signer) Option { return func(c *config) { c.signer = s } }

// WithSubaccount sets the subaccount id used by private REST and WS calls.
func WithSubaccount(id int64) Option { return func(c *config) { c.subaccount = id } }

// WithConnectWS controls whether [NewClient] dials the WebSocket up front.
//
// When true, NewClient calls Connect (and Login if a signer is configured)
// before returning. Default is false — most callers prefer to call
// c.WS.Connect(ctx) and c.WS.Login(ctx) themselves so they can plumb their
// own context and surface the resulting error.
func WithConnectWS(b bool) Option { return func(c *config) { c.connectWS = b } }

// WithInstrumentPreload kicks off a background fetch of every live
// instrument via the REST client immediately after construction.
//
// The preload populates the shared instrument metadata cache used by
// signed actions on both transports. See
// [rest.WithInstrumentPreload] for details.
func WithInstrumentPreload() Option { return func(c *config) { c.preloadInsts = true } }

// WithUserAgent overrides the User-Agent header on both REST requests
// and the WS upgrade. Default is "go-derive/<version>". See
// [rest.WithUserAgent] and [ws.WithUserAgent].
func WithUserAgent(ua string) Option { return func(c *config) { c.userAgent = ua } }

// WithRateLimit configures the per-client token-bucket rate limiter
// for both REST and WS. tps is the sustained transactions-per-second;
// burst sets the bucket capacity (capacity = tps × burst). Pass tps
// <= 0 to disable. Default: 10 TPS, burst 5. See
// [rest.WithRateLimit] and [ws.WithRateLimit].
func WithRateLimit(tps, burst float64) Option {
	return func(c *config) {
		c.tps = &tps
		c.burst = &burst
	}
}

// WithSignatureExpiry sets the seconds-from-now expiry on signed actions
// for both REST and WS. Default is 300 (5 minutes). See
// [rest.WithSignatureExpiry] and [ws.WithSignatureExpiry].
func WithSignatureExpiry(seconds int64) Option {
	return func(c *config) { c.expiry = &seconds }
}

// WithHTTPClient swaps in a custom *http.Client for the REST transport
// (for custom transports, proxies, mocking). REST-only — the WS
// transport doesn't use *http.Client for its persistent connection.
// See [rest.WithHTTPClient].
func WithHTTPClient(h *http.Client) Option { return func(c *config) { c.httpClient = h } }

// WithHTTPTimeout sets the total per-request timeout on the REST
// client's default *http.Client. Default is 30 seconds. Pass 0 to
// disable. Ignored when [WithHTTPClient] is also set. REST-only.
// See [rest.WithHTTPTimeout].
func WithHTTPTimeout(d time.Duration) Option {
	return func(c *config) { c.httpTimeout = &d }
}

// WithPingInterval overrides the periodic application-level ping
// cadence on the WS transport. Default is 20 seconds. WS-only.
// See [ws.WithPingInterval].
func WithPingInterval(d time.Duration) Option {
	return func(c *config) { c.pingInterval = &d }
}

// WithReconnect controls whether the WS client auto-reconnects after
// the underlying connection drops. Default is true. WS-only.
// See [ws.WithReconnect].
func WithReconnect(enabled bool) Option {
	return func(c *config) { c.reconnect = &enabled }
}

// WithOnReconnect installs a callback invoked once per WS reconnect
// cycle, after the redial, internal re-login, and resubscribe all
// complete. err is nil on full recovery; non-nil when the post-redial
// chain partially failed. WS-only. See [ws.WithOnReconnect] for full
// semantics and the safety contract.
func WithOnReconnect(fn func(err error)) Option {
	return func(c *config) { c.onReconnect = fn }
}

// NewClient constructs a [Client] from the supplied options.
//
// One of [WithMainnet], [WithTestnet], or [WithCustomNetwork] must be
// supplied. If both REST and WS construction succeed but a subsequent
// optional Connect/Login fails, NewClient closes both clients before
// returning the error.
func NewClient(opts ...Option) (*Client, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.network.Network == netconf.NetworkUnknown {
		return nil, derrors.ErrInvalidConfig
	}

	restOpts := []rest.Option{
		rest.WithCustomNetwork(cfg.network),
		rest.WithSigner(cfg.signer),
		rest.WithSubaccount(cfg.subaccount),
	}
	wsOpts := []ws.Option{
		ws.WithCustomNetwork(cfg.network),
		ws.WithSigner(cfg.signer),
		ws.WithSubaccount(cfg.subaccount),
	}

	// Common options forwarded to both transports.
	if cfg.userAgent != "" {
		restOpts = append(restOpts, rest.WithUserAgent(cfg.userAgent))
		wsOpts = append(wsOpts, ws.WithUserAgent(cfg.userAgent))
	}
	if cfg.tps != nil && cfg.burst != nil {
		restOpts = append(restOpts, rest.WithRateLimit(*cfg.tps, *cfg.burst))
		wsOpts = append(wsOpts, ws.WithRateLimit(*cfg.tps, *cfg.burst))
	}
	if cfg.expiry != nil {
		restOpts = append(restOpts, rest.WithSignatureExpiry(*cfg.expiry))
		wsOpts = append(wsOpts, ws.WithSignatureExpiry(*cfg.expiry))
	}
	if cfg.preloadInsts {
		restOpts = append(restOpts, rest.WithInstrumentPreload())
		wsOpts = append(wsOpts, ws.WithInstrumentPreload())
	}

	// REST-only options.
	if cfg.httpClient != nil {
		restOpts = append(restOpts, rest.WithHTTPClient(cfg.httpClient))
	}
	if cfg.httpTimeout != nil {
		restOpts = append(restOpts, rest.WithHTTPTimeout(*cfg.httpTimeout))
	}

	// WS-only options.
	if cfg.pingInterval != nil {
		wsOpts = append(wsOpts, ws.WithPingInterval(*cfg.pingInterval))
	}
	if cfg.reconnect != nil {
		wsOpts = append(wsOpts, ws.WithReconnect(*cfg.reconnect))
	}
	if cfg.onReconnect != nil {
		wsOpts = append(wsOpts, ws.WithOnReconnect(cfg.onReconnect))
	}

	r, err := rest.New(restOpts...)
	if err != nil {
		return nil, err
	}
	w, err := ws.New(wsOpts...)
	if err != nil {
		_ = r.Close()
		return nil, err
	}

	c := &Client{REST: r, WS: w, cfg: cfg.network}
	if cfg.connectWS {
		ctx := context.Background()
		if err := w.Connect(ctx); err != nil {
			_ = c.Close()
			return nil, err
		}
		if cfg.signer != nil {
			if err := w.Login(ctx); err != nil {
				_ = c.Close()
				return nil, err
			}
		}
	}
	return c, nil
}

// Close releases both transports' resources. Idempotent.
func (c *Client) Close() error {
	wsErr := c.WS.Close()
	restErr := c.REST.Close()
	if wsErr != nil {
		return wsErr
	}
	return restErr
}

// Network returns the active network configuration. Useful for diagnostics
// and for plumbing the same config into auxiliary tooling.
func (c *Client) Network() netconf.Config { return c.cfg }
