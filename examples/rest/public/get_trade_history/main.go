// Paginates the public trade tape — page-sized 5, prints the latest 5
// trades plus the total record count and page count Derive reports.
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
	instrument := os.Getenv("DERIVE_INSTRUMENT")
	if instrument == "" {
		instrument = "BTC-PERP"
	}

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

	trades, page, err := c.GetPublicTradeHistory(ctx,
		types.PublicTradeHistoryQuery{InstrumentName: instrument},
		types.PageRequest{PageSize: 5})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "trades returned:", len(trades))
	fmt.Printf("%-30s %v\n", "total records:", page.Count)
	fmt.Printf("%-30s %v\n", "total pages:", page.NumPages)
	for _, t := range trades {
		fmt.Printf("%-30s %v\n", string(t.Direction)+" "+t.TradeAmount.String()+":", t.TradePrice)
	}
}
