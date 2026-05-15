// Fetches the per-expiry settlement prices for one currency's
// option market. Pre-settlement entries return Price as zero.
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
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		currency = "BTC"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	prices, err := c.GetOptionSettlementPrices(ctx, types.OptionSettlementPricesQuery{Currency: currency})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "expiry count:", len(prices))
	for i, p := range prices {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "expiry:", p.ExpiryDate)
		fmt.Printf("%-30s %v\n", "  utc_expiry_sec:", p.UTCExpirySec)
		fmt.Printf("%-30s %v\n", "  price:", p.Price.String())
	}
}
