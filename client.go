// Top-level convenience facade for the SDK. [Client] composes a
// *[RestClient] and a *[WsClient] sharing the same signer, subaccount
// and network configuration.
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
//	insts, _ := c.REST.GetInstruments(ctx, "BTC", derive.InstrumentTypePerp)
//	c.WS.Connect(ctx)
//	c.WS.Login(ctx)
//
// # When to skip the facade
//
// Use [NewRestClient] or [NewWsClient] directly when you only need one
// transport — both expose the full RPC method surface independently.
// The facade is just a shortcut for the common case where you want
// both.

package derive

import (
	"context"
)

// Client bundles a REST client and a WebSocket client sharing the same
// signer, subaccount and network configuration.
//
// Both REST and WS are exported so callers can reach the underlying
// client for transport-specific options (e.g. c.WS.Connect(ctx)). The
// zero value is not usable; obtain one from [NewClient].
type Client struct {
	// REST is the HTTP-backed JSON-RPC client.
	REST *RestClient
	// WS is the WebSocket-backed JSON-RPC + subscription client.
	WS *WsClient

	cfg NetworkConfig
}

// WithConnectWS controls whether [NewClient] dials the WebSocket up
// front.
//
// When true, NewClient calls Connect (and Login if a signer is
// configured) before returning. Default is false — most callers prefer
// to call c.WS.Connect(ctx) and c.WS.Login(ctx) themselves so they can
// plumb their own context and surface the resulting error.
//
// Facade-only — ignored by [NewRestClient] and [NewWsClient].
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
	if cfg.network.Network == NetworkUnknown {
		return nil, ErrInvalidConfig
	}

	r, err := NewRestClient(
		WithCustomNetwork(cfg.network),
		WithSigner(cfg.signer),
		WithSubaccount(cfg.subaccount),
	)
	if err != nil {
		return nil, err
	}
	w, err := NewWsClient(
		WithCustomNetwork(cfg.network),
		WithSigner(cfg.signer),
		WithSubaccount(cfg.subaccount),
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

// Network returns the active network configuration. Useful for
// diagnostics and for plumbing the same config into auxiliary tooling.
func (c *Client) Network() NetworkConfig { return c.cfg }
