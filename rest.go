// HTTP-backed JSON-RPC client lives in this file: [RestClient] and its
// [Option] knobs.
//
// # Concurrency
//
// [RestClient] is safe for concurrent use. Construct one per process and
// share it across goroutines; every method is goroutine-safe.
//
// # Authentication
//
// Pass a [Signer] via [WithSigner] to enable private endpoints. The
// transport adds X-LyraWallet, X-LyraTimestamp and X-LyraSignature
// headers to every request automatically. Public endpoints work without
// a signer.
//
// # Method surface
//
// [RestClient] embeds *[API] (still exported until commit 13 brings the
// facade into root), exposing every documented JSON-RPC method as a
// regular Go method:
//
//	c, _ := derive.NewRestClient(derive.WithMainnet(), derive.WithSigner(s), derive.WithSubaccount(123))
//	instruments, _ := c.GetInstruments(ctx, "BTC", derive.InstrumentTypePerp)
//	c.PlaceOrder(ctx, derive.PlaceOrderInput{...})
//
// # Errors
//
// REST errors arrive as [*APIError] when the server rejected the call,
// or [*ConnectionError] for transport-level failures. Both compose with
// errors.Is and errors.As.

package derive

import (
	"context"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/internal/transport"
)

// RestClient is the HTTP JSON-RPC client.
//
// Construct one with [NewRestClient], plus the desired With* options.
// The zero value is not usable.
type RestClient struct {
	*API
	http   *transport.HTTPTransport
	signer Signer
	cfg    NetworkConfig
}

// NewRestClient constructs a [RestClient] with the given options.
//
// One of [WithMainnet], [WithTestnet], or [WithCustomNetwork] must be
// supplied; without a network selection NewRestClient returns
// [ErrInvalidConfig].
//
// All other options have sensible defaults: 30-second HTTP timeout,
// 10 TPS rate limit with 5x burst, 5-minute signature expiry, and the
// SDK's default User-Agent. See the With* helpers for what's tunable.
func NewRestClient(opts ...Option) (*RestClient, error) {
	c := &config{
		userAgent: UserAgent(),
		expiry:    300,
		tps:       10,
		burst:     5,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.network.Network == NetworkUnknown {
		return nil, ErrInvalidConfig
	}

	hdrProv := func(ctx context.Context, _ string, _ []byte) (http.Header, error) {
		if c.signer == nil {
			return nil, nil
		}
		return HTTPHeaders(ctx, c.signer, time.Now())
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

	api := &API{
		T:               httpT,
		Signer:          c.signer,
		Domain:          c.network.EIP712Domain(),
		Subaccount:      c.subaccount,
		Nonces:          NewNonceGen(),
		SignatureExpiry: c.expiry,
	}
	api.SetTradeModule(common.HexToAddress(c.network.Contracts.TradeModule))

	return &RestClient{API: api, http: httpT, signer: c.signer, cfg: c.network}, nil
}

// Close releases transport-level resources. The [RestClient] is unusable
// afterwards. Close is idempotent.
func (c *RestClient) Close() error { return c.http.Close() }

// Network returns the active network configuration. Useful for
// diagnostics and for plumbing the same configuration into a separate
// WS client.
func (c *RestClient) Network() NetworkConfig { return c.cfg }

// Option configures a [RestClient] at construction time. Compose any
// number of With* helpers and pass them to [NewRestClient].
//
// Commit 12 widens this type to also configure [WsClient] and the
// top-level facade; the WS-only and REST-only knobs are documented per
// helper.
type Option func(*config)

type config struct {
	network    NetworkConfig
	signer     Signer
	subaccount int64
	// REST-only.
	httpClient *http.Client
	userAgent  string
	tps        float64
	burst      float64
	expiry     int64
	// WS-only.
	pingInterval time.Duration
	reconnect    bool
}

// WithMainnet selects Derive's mainnet endpoints (chain id 957). Required
// unless WithTestnet or WithCustomNetwork is used.
func WithMainnet() Option { return func(c *config) { c.network = Mainnet() } }

// WithTestnet selects Derive's demo (testnet) endpoints (chain id 901).
// Use this for development before promoting to mainnet.
func WithTestnet() Option { return func(c *config) { c.network = Testnet() } }

// WithCustomNetwork overrides the entire network configuration. Use it
// for staging or vendored deployments where the default mainnet/testnet
// endpoints do not apply.
func WithCustomNetwork(cfg NetworkConfig) Option { return func(c *config) { c.network = cfg } }

// WithSigner attaches an auth [Signer] used for both REST auth headers
// and per-action EIP-712 signing.
//
// Without a signer, only public endpoints will succeed; private
// endpoints return [ErrUnauthorized].
func WithSigner(s Signer) Option { return func(c *config) { c.signer = s } }

// WithSubaccount sets the subaccount id used by private methods that
// don't take an explicit subaccount.
func WithSubaccount(id int64) Option { return func(c *config) { c.subaccount = id } }

// WithHTTPClient swaps in a custom *http.Client (for custom transports,
// proxies, mocking, etc.). The default is a *http.Client with a
// 30-second timeout.
//
// REST-only — ignored by [WsClient].
func WithHTTPClient(h *http.Client) Option { return func(c *config) { c.httpClient = h } }

// WithUserAgent overrides the default User-Agent header (which is
// "go-derive/<version>"). Useful for distinguishing your fleet in
// Derive-side logs.
func WithUserAgent(ua string) Option { return func(c *config) { c.userAgent = ua } }

// WithRateLimit configures the per-client token-bucket rate limiter.
//
// tps is the sustained rate in transactions per second; burst is a
// multiplier giving the bucket capacity (so capacity = tps × burst).
// Pass tps <= 0 to disable limiting entirely. Defaults: 10 TPS, burst 5.
func WithRateLimit(tps, burst float64) Option {
	return func(c *config) {
		c.tps = tps
		c.burst = burst
	}
}

// WithSignatureExpiry sets the seconds-from-now expiry on signed
// actions.
//
// The default is 300 (5 minutes). Use a longer value if you sign orders
// far in advance of submission, or a shorter value if you need tighter
// replay-protection bounds.
func WithSignatureExpiry(seconds int64) Option {
	return func(c *config) { c.expiry = seconds }
}
