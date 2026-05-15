// Returns the invite code allocated to one wallet plus its
// remaining-uses counter. Optional env: DERIVE_WALLET.
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

	res, err := c.GetInviteCode(ctx, types.InviteCodeQuery{Wallet: os.Getenv("DERIVE_WALLET")})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "code:", res.Code)
	fmt.Printf("%-30s %v\n", "remaining_uses (-1=unlimited):", res.RemainingUses)
}
