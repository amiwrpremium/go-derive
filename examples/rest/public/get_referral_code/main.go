// Returns the referral code currently associated with one wallet.
// Optional env: DERIVE_WALLET (defaults to signer's wallet if any).
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork)
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	code, err := c.GetReferralCode(ctx, types.ReferralCodeQuery{Wallet: os.Getenv("DERIVE_WALLET")})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "referral_code:", code)
}
