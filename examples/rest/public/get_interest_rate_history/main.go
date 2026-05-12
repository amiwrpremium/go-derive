// Fetches historical USDC borrow / supply APY prints over the last
// 24 hours.
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

	end := time.Now().Unix()
	start := end - 24*3600

	rates, page, err := c.GetInterestRateHistory(ctx, types.InterestRateHistoryQuery{
		FromSec: start,
		ToSec:   end,
	}, types.PageRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "rate prints:", len(rates))
	fmt.Printf("%-30s %v\n", "total pages:", page.NumPages)
	for i, r := range rates {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "at sec:", r.TimestampSec)
		fmt.Printf("%-30s %v\n", "  borrow / supply APY:", r.BorrowAPY.String()+" / "+r.SupplyAPY.String())
	}
}
