// Package rest is the HTTP-backed client for Derive's JSON-RPC API.
//
// # Concurrency
//
// [Client] is safe for concurrent use. Construct one per process and share
// it across goroutines; every method is goroutine-safe.
//
// # Authentication
//
// Pass a [github.com/amiwrpremium/go-derive/pkg/auth.Signer] via [WithSigner]
// to enable private endpoints. The transport adds X-LyraWallet,
// X-LyraTimestamp and X-LyraSignature headers to every request automatically.
// Public endpoints work without a signer.
//
// # Method surface
//
// [Client] embeds *methods.API (an internal type), exposing every documented
// JSON-RPC method as a regular Go method:
//
//	c, _ := rest.New(rest.WithMainnet(), rest.WithSigner(s), rest.WithSubaccount(123))
//	instruments, _ := c.GetInstruments(ctx, "BTC", enums.InstrumentTypePerp)
//	c.PlaceOrder(ctx, types.PlaceOrderInput{...})
//
// # Errors
//
// REST errors arrive as [*github.com/amiwrpremium/go-derive/pkg/errors.APIError]
// when the server rejected the call, or [*github.com/amiwrpremium/go-derive/pkg/errors.ConnectionError]
// for transport-level failures. Both compose with errors.Is and errors.As.
package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"

	goderive "github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/methods"
	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/internal/transport"
	"github.com/amiwrpremium/go-derive/pkg/auth"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

// Client is the HTTP JSON-RPC client.
//
// Construct one with [New], plus the desired With* options. The zero value
// is not usable.
type Client struct {
	*methods.API
	http   *transport.HTTPTransport
	signer auth.Signer
	cfg    netconf.Config
}

// New constructs a [Client] with the given options.
//
// One of [WithMainnet], [WithTestnet], or [WithCustomNetwork] must be supplied;
// without a network selection New returns
// [github.com/amiwrpremium/go-derive/pkg/errors.ErrInvalidConfig].
//
// All other options have sensible defaults: 30-second HTTP timeout, 10 TPS
// rate limit with 5x burst, 5-minute signature expiry, and the SDK's default
// User-Agent. See the With* helpers in options.go for what's tunable.
func New(opts ...Option) (*Client, error) {
	c := &config{
		userAgent: goderive.UserAgent(),
		expiry:    300, // 5 minutes
		tps:       10,
		burst:     5,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.network.Network == netconf.NetworkUnknown {
		return nil, derrors.ErrInvalidConfig
	}

	hdrProv := func(ctx context.Context, _ string, _ []byte) (http.Header, error) {
		if c.signer == nil {
			return nil, nil
		}
		return auth.HTTPHeaders(ctx, c.signer, time.Now())
	}

	httpT, err := transport.NewHTTP(transport.HTTPConfig{
		URL:       c.network.HTTPURL,
		Client:    c.httpClient,
		UserAgent: c.userAgent,
		Limiter:   transport.NewRateLimiter(c.tps, c.burst),
		Headers:   hdrProv,
	})
	if err != nil {
		return nil, err
	}

	api := &methods.API{
		T:               httpT,
		Signer:          c.signer,
		Domain:          c.network.EIP712Domain(),
		Subaccount:      c.subaccount,
		Nonces:          auth.NewNonceGen(),
		SignatureExpiry: c.expiry,
	}
	api.SetTradeModule(common.HexToAddress(c.network.Contracts.TradeModule))
	api.SetRFQModule(common.HexToAddress(c.network.Contracts.RFQModule))

	client := &Client{API: api, http: httpT, signer: c.signer, cfg: c.network}

	if c.preloadAllInsts {
		go func() {
			_ = api.PreloadAllInstruments(context.Background())
		}()
	}
	return client, nil
}

// Close releases transport-level resources. The [Client] is unusable
// afterwards. Close is idempotent.
func (c *Client) Close() error { return c.http.Close() }

// Network returns the active network configuration. Useful for diagnostics
// and for plumbing the same configuration into a separate WS client.
func (c *Client) Network() netconf.Config { return c.cfg }
