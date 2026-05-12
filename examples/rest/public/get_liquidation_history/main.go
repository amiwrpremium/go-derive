// Lists the network-wide liquidation auction history (paginated).
// Counterpart to private/get_liquidation_history (single-subaccount).
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

	auctions, page, err := c.GetPublicLiquidationHistory(ctx, types.LiquidationHistoryQuery{}, types.PageRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "auctions:", len(auctions))
	fmt.Printf("%-30s %v\n", "total pages:", page.NumPages)
	for i, a := range auctions {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "auction:", a.AuctionID)
		fmt.Printf("%-30s %v\n", "  type:", string(a.AuctionType))
		fmt.Printf("%-30s %v\n", "  subaccount:", a.SubaccountID)
	}
}
