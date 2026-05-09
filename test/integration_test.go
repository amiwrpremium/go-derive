//go:build integration

// Package integration_test contains live-network integration tests for
// go-derive.
//
// All tests in this directory are gated by the `integration` build tag, so
// the default `go test ./...` is unaffected. Run them with:
//
//	go test -tags=integration ./test/...
//
// See test/README.md for the env vars each test needs and the safety rails
// around live order placement.
package integration_test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	goderive "github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/derive"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

// integrationEnv is the configuration loaded from environment variables.
type integrationEnv struct {
	network    goderive.NetworkConfig
	sessionKey string
	owner      common.Address
	subaccount int64
	instrument string
	baseAsset  common.Address
	liveOrders bool
}

// loadEnv reads the test configuration from environment variables.
//
// The defaults target Derive testnet with BTC-PERP. Tests skip if a
// required field is missing.
func loadEnv(t *testing.T) integrationEnv {
	t.Helper()
	env := integrationEnv{
		instrument: getenv("DERIVE_INSTRUMENT", "BTC-PERP"),
	}

	switch getenv("DERIVE_NETWORK", "testnet") {
	case "mainnet":
		env.network = goderive.Mainnet()
	default:
		env.network = goderive.Testnet()
	}

	env.sessionKey = os.Getenv("DERIVE_SESSION_KEY")
	if owner := os.Getenv("DERIVE_OWNER"); owner != "" {
		env.owner = common.HexToAddress(owner)
	}
	if sub := os.Getenv("DERIVE_SUBACCOUNT"); sub != "" {
		n, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			t.Fatalf("DERIVE_SUBACCOUNT=%q: %v", sub, err)
		}
		env.subaccount = n
	}
	if asset := os.Getenv("DERIVE_BASE_ASSET"); asset != "" {
		env.baseAsset = common.HexToAddress(asset)
	}
	env.liveOrders = os.Getenv("DERIVE_RUN_LIVE_ORDERS") == "1"
	return env
}

// hasPrivateCreds reports whether the env carries enough to authenticate.
func (e integrationEnv) hasPrivateCreds() bool {
	return e.sessionKey != "" && e.subaccount != 0
}

// requireSigner skips the test when private creds are absent; otherwise
// returns a configured Signer.
func requireSigner(t *testing.T, env integrationEnv) goderive.Signer {
	t.Helper()
	if !env.hasPrivateCreds() {
		t.Skip("private creds not configured (DERIVE_SESSION_KEY + DERIVE_SUBACCOUNT)")
	}
	var (
		s   goderive.Signer
		err error
	)
	if env.owner == (common.Address{}) {
		s, err = goderive.NewLocalSigner(env.sessionKey)
	} else {
		s, err = goderive.NewSessionKeySigner(env.sessionKey, env.owner)
	}
	if err != nil {
		t.Fatalf("build signer: %v", err)
	}
	return s
}

// newRESTClient wires up a public-only REST client.
func newRESTClient(t *testing.T, env integrationEnv) *goderive.RestClient {
	t.Helper()
	c, err := goderive.NewRestClient(goderive.WithCustomNetwork(env.network))
	if err != nil {
		t.Fatalf("derive.NewRestClient: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

// newAuthRESTClient wires up an authenticated REST client. Skips if creds
// are missing.
func newAuthRESTClient(t *testing.T, env integrationEnv) *goderive.RestClient {
	t.Helper()
	signer := requireSigner(t, env)
	c, err := goderive.NewRestClient(
		goderive.WithCustomNetwork(env.network),
		goderive.WithSigner(signer),
		goderive.WithSubaccount(env.subaccount),
	)
	if err != nil {
		t.Fatalf("derive.NewRestClient: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

// newWSClient wires up a public-only WS client and connects.
func newWSClient(t *testing.T, env integrationEnv) *ws.Client {
	t.Helper()
	c, err := ws.New(ws.WithCustomNetwork(env.network))
	if err != nil {
		t.Fatalf("ws.New: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := c.Connect(ctx); err != nil {
		_ = c.Close()
		t.Fatalf("ws.Connect: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

// newAuthWSClient wires up an authenticated WS client, connects and logs
// in. Skips if creds are missing.
func newAuthWSClient(t *testing.T, env integrationEnv) *ws.Client {
	t.Helper()
	signer := requireSigner(t, env)
	c, err := ws.New(
		ws.WithCustomNetwork(env.network),
		ws.WithSigner(signer),
		ws.WithSubaccount(env.subaccount),
	)
	if err != nil {
		t.Fatalf("ws.New: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := c.Connect(ctx); err != nil {
		_ = c.Close()
		t.Fatalf("ws.Connect: %v", err)
	}
	if err := c.Login(ctx); err != nil {
		_ = c.Close()
		t.Fatalf("ws.Login: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

// newDeriveClient wires up the top-level facade.
func newDeriveClient(t *testing.T, env integrationEnv) *derive.Client {
	t.Helper()
	opts := []derive.Option{derive.WithCustomNetwork(env.network)}
	if env.hasPrivateCreds() {
		opts = append(opts, derive.WithSigner(requireSigner(t, env)),
			derive.WithSubaccount(env.subaccount))
	}
	c, err := derive.NewClient(opts...)
	if err != nil {
		t.Fatalf("derive.NewClient: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

// withTimeout returns a 30-second context for a single integration call.
func withTimeout(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
