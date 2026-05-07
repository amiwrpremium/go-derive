// Package rest — see client.go for the overview.
package rest

import (
	"net/http"

	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

// Option configures a [Client] at construction time. Compose any number of
// With* helpers and pass them to [New].
type Option func(*config)

type config struct {
	network    netconf.Config
	signer     auth.Signer
	subaccount int64
	httpClient *http.Client
	userAgent  string
	tps        float64
	burst      float64
	expiry     int64
}

// WithMainnet selects Derive's mainnet endpoints (chain id 957). Required
// unless WithTestnet or WithCustomNetwork is used.
func WithMainnet() Option { return func(c *config) { c.network = netconf.Mainnet() } }

// WithTestnet selects Derive's demo (testnet) endpoints (chain id 901).
// Use this for development before promoting to mainnet.
func WithTestnet() Option { return func(c *config) { c.network = netconf.Testnet() } }

// WithCustomNetwork overrides the entire network configuration. Use it for
// staging or vendored deployments where the default mainnet/testnet
// endpoints do not apply.
func WithCustomNetwork(cfg netconf.Config) Option { return func(c *config) { c.network = cfg } }

// WithSigner attaches an auth [github.com/amiwrpremium/go-derive/pkg/auth.Signer]
// used for both REST auth headers and per-action EIP-712 signing.
//
// Without a signer, only public endpoints will succeed; private endpoints
// return [github.com/amiwrpremium/go-derive/pkg/errors.ErrUnauthorized].
func WithSigner(s auth.Signer) Option { return func(c *config) { c.signer = s } }

// WithSubaccount sets the subaccount id used by private methods that don't
// take an explicit subaccount.
func WithSubaccount(id int64) Option { return func(c *config) { c.subaccount = id } }

// WithHTTPClient swaps in a custom *http.Client (for custom transports,
// proxies, mocking, etc.). The default is a *http.Client with a 30-second
// timeout.
func WithHTTPClient(h *http.Client) Option { return func(c *config) { c.httpClient = h } }

// WithUserAgent overrides the default User-Agent header (which is
// "go-derive/<version>"). Useful for distinguishing your fleet in
// Derive-side logs.
func WithUserAgent(ua string) Option { return func(c *config) { c.userAgent = ua } }

// WithRateLimit configures the per-client token-bucket rate limiter.
//
// tps is the sustained rate in transactions per second; burst is a
// multiplier giving the bucket capacity (so capacity = tps × burst). Pass
// tps <= 0 to disable limiting entirely. Defaults: 10 TPS, burst 5.
func WithRateLimit(tps, burst float64) Option {
	return func(c *config) {
		c.tps = tps
		c.burst = burst
	}
}

// WithSignatureExpiry sets the seconds-from-now expiry on signed actions.
//
// The default is 300 (5 minutes). Use a longer value if you sign orders
// far in advance of submission, or a shorter value if you need tighter
// replay-protection bounds.
func WithSignatureExpiry(seconds int64) Option {
	return func(c *config) { c.expiry = seconds }
}
