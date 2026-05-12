// Fetches per-asset margin parameters and protocol-asset addresses
// for one underlying currency.
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
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		currency = "ETH"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := c.GetCurrency(ctx, currency)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "currency:", cur.Currency)
	fmt.Printf("%-30s %v\n", "spot_price:", cur.SpotPrice.String())
	fmt.Printf("%-30s %v\n", "market_type:", cur.MarketType)
	fmt.Printf("%-30s %v\n", "manager count:", len(cur.Managers))
	fmt.Printf("%-30s %v\n", "instrument_types:", cur.InstrumentTypes)
}
