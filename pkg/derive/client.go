// Package derive is the top-level convenience facade for the SDK. It
// composes a *[github.com/amiwrpremium/go-derive.RestClient] and a
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
//	insts, _ := c.REST.GetInstruments(ctx, "BTC", enums.InstrumentTypePerp)
//	c.WS.Connect(ctx)
//	c.WS.Login(ctx)
//
// # When to skip the facade
//
// Use the root *[github.com/amiwrpremium/go-derive.RestClient] or
// pkg/ws directly when you only need one transport — both expose the
// full RPC method surface independently. The facade is just a shortcut
// for the common case where you want both.
package derive

import (
	"context"

	"github.com/amiwrpremium/go-derive"
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
	REST *derive.RestClient
	// WS is the WebSocket-backed JSON-RPC + subscription client.
	WS *ws.Client

	cfg derive.NetworkConfig
}

// Option configures a [Client] at construction time. Compose any number of
// With* helpers and pass them to [NewClient].
type Option func(*config)

type config struct {
	network    derive.NetworkConfig
	signer     derive.Signer
	subaccount int64
	connectWS  bool
}

// WithMainnet selects Derive's mainnet endpoints.
func WithMainnet() Option { return func(c *config) { c.network = derive.Mainnet() } }

// WithTestnet selects Derive's demo (testnet) endpoints.
func WithTestnet() Option { return func(c *config) { c.network = derive.Testnet() } }

// WithCustomNetwork overrides the entire network configuration. Use it for
// staging or vendored deployments.
func WithCustomNetwork(cfg derive.NetworkConfig) Option { return func(c *config) { c.network = cfg } }

// WithSigner attaches an auth [github.com/amiwrpremium/go-derive/pkg/derive.Signer].
// Without one, only public endpoints work.
func WithSigner(s derive.Signer) Option { return func(c *config) { c.signer = s } }

// WithSubaccount sets the subaccount id used by private REST and WS calls.
func WithSubaccount(id int64) Option { return func(c *config) { c.subaccount = id } }

// WithConnectWS controls whether [NewClient] dials the WebSocket up front.
//
// When true, NewClient calls Connect (and Login if a signer is configured)
// before returning. Default is false — most callers prefer to call
// c.WS.Connect(ctx) and c.WS.Login(ctx) themselves so they can plumb their
// own context and surface the resulting error.
func WithConnectWS(b bool) Option { return func(c *config) { c.connectWS = b } }

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
	if cfg.network.Network == derive.NetworkUnknown {
		return nil, derive.ErrInvalidConfig
	}

	r, err := derive.NewRestClient(
		derive.WithCustomNetwork(cfg.network),
		derive.WithSigner(cfg.signer),
		derive.WithSubaccount(cfg.subaccount),
	)
	if err != nil {
		return nil, err
	}
	w, err := ws.New(
		ws.WithCustomNetwork(cfg.network),
		ws.WithSigner(cfg.signer),
		ws.WithSubaccount(cfg.subaccount),
	)
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
func (c *Client) Network() derive.NetworkConfig { return c.cfg }
