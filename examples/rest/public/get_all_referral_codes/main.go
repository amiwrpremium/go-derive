// Lists every valid referral code for the configured signer's
// wallet (or omits the wallet param when no signer is configured).
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
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

	codes, err := c.GetAllReferralCodes(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "codes:", len(codes))
	for i, r := range codes {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "code:", r.ReferralCode)
		fmt.Printf("%-30s %v\n", "  wallet:", r.Wallet)
	}
}
