// Package example holds shared setup helpers for the runnable example
// programs under examples/. Every example main.go uses these so the
// per-example file stays focused on the API call being demonstrated.
//
// Environment variables (all optional except as noted):
//
//	DERIVE_NETWORK     mainnet|testnet (default testnet)
//	DERIVE_INSTRUMENT  e.g. BTC-PERP   (default BTC-PERP)
//	DERIVE_BASE_ASSET  on-chain asset address used for trade-module signing
//	DERIVE_SESSION_KEY hex private key (required for private examples)
//	DERIVE_OWNER       owner address    (required for SessionKeySigner)
//	DERIVE_SUBACCOUNT  numeric id       (required for private examples)
package example

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	goderive "github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/derive"
	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

// Network returns the configured Derive network (default: testnet).
func Network() goderive.NetworkConfig {
	switch os.Getenv("DERIVE_NETWORK") {
	case "mainnet":
		return goderive.Mainnet()
	default:
		return goderive.Testnet()
	}
}

// Instrument returns DERIVE_INSTRUMENT or "BTC-PERP".
func Instrument() string {
	if v := os.Getenv("DERIVE_INSTRUMENT"); v != "" {
		return v
	}
	return "BTC-PERP"
}

// BaseAsset returns DERIVE_BASE_ASSET as a common.Address; zero if unset.
func BaseAsset() common.Address {
	if v := os.Getenv("DERIVE_BASE_ASSET"); v != "" {
		return common.HexToAddress(v)
	}
	return common.Address{}
}

// Subaccount returns DERIVE_SUBACCOUNT parsed as int64; logs.Fatal if
// missing. Use for examples that require a private endpoint.
func Subaccount() int64 {
	v := os.Getenv("DERIVE_SUBACCOUNT")
	if v == "" {
		log.Fatal("DERIVE_SUBACCOUNT required for this example")
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_SUBACCOUNT=%q: %v", v, err)
	}
	return n
}

// MustSigner returns an goderive.Signer built from the env, or fatal-exits
// when DERIVE_SESSION_KEY is missing.
func MustSigner() goderive.Signer {
	key := os.Getenv("DERIVE_SESSION_KEY")
	if key == "" {
		log.Fatal("DERIVE_SESSION_KEY required for this example")
	}
	if owner := os.Getenv("DERIVE_OWNER"); owner != "" {
		s, err := goderive.NewSessionKeySigner(key, common.HexToAddress(owner))
		if err != nil {
			log.Fatalf("session key signer: %v", err)
		}
		return s
	}
	s, err := goderive.NewLocalSigner(key)
	if err != nil {
		log.Fatalf("local signer: %v", err)
	}
	return s
}

// MustRESTPublic returns a public-only REST client. Always succeeds.
func MustRESTPublic() *rest.Client {
	c, err := rest.New(rest.WithCustomNetwork(Network()))
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	return c
}

// MustRESTPrivate returns an authenticated REST client using env creds.
func MustRESTPrivate() *rest.Client {
	c, err := rest.New(
		rest.WithCustomNetwork(Network()),
		rest.WithSigner(MustSigner()),
		rest.WithSubaccount(Subaccount()),
	)
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	return c
}

// MustWSPublic returns a connected public WS client. Caller must defer Close.
func MustWSPublic(ctx context.Context) *ws.Client {
	c, err := ws.New(ws.WithCustomNetwork(Network()))
	if err != nil {
		log.Fatalf("ws.New: %v", err)
	}
	if err := c.Connect(ctx); err != nil {
		_ = c.Close()
		log.Fatalf("ws.Connect: %v", err)
	}
	return c
}

// MustWSPrivate returns a connected, logged-in WS client.
func MustWSPrivate(ctx context.Context) *ws.Client {
	c, err := ws.New(
		ws.WithCustomNetwork(Network()),
		ws.WithSigner(MustSigner()),
		ws.WithSubaccount(Subaccount()),
	)
	if err != nil {
		log.Fatalf("ws.New: %v", err)
	}
	if err := c.Connect(ctx); err != nil {
		_ = c.Close()
		log.Fatalf("ws.Connect: %v", err)
	}
	if err := c.Login(ctx); err != nil {
		_ = c.Close()
		log.Fatalf("ws.Login: %v", err)
	}
	return c
}

// MustDerivePublic returns the top-level facade with public-only access.
func MustDerivePublic() *derive.Client {
	c, err := derive.NewClient(derive.WithCustomNetwork(Network()))
	if err != nil {
		log.Fatalf("derive.NewClient: %v", err)
	}
	return c
}

// MustDerivePrivate returns the top-level facade with creds.
func MustDerivePrivate() *derive.Client {
	c, err := derive.NewClient(
		derive.WithCustomNetwork(Network()),
		derive.WithSigner(MustSigner()),
		derive.WithSubaccount(Subaccount()),
	)
	if err != nil {
		log.Fatalf("derive.NewClient: %v", err)
	}
	return c
}

// Timeout returns a 30-second context for one-shot calls.
func Timeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// LongTimeout returns a 2-minute context for subscription examples.
func LongTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Minute)
}

// Fatal is a thin wrapper that logs and exits, kept short for one-line
// `if err != nil { example.Fatal(err) }` use at example call sites.
func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Print formats a result using the same convention every example uses,
// keeping their output uniform.
func Print(label string, v any) {
	fmt.Printf("%-30s %v\n", label+":", v)
}
