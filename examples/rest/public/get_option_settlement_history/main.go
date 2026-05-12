// Lists platform-wide option settlements. Public — no auth needed.
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

	settlements, page, err := c.GetPublicOptionSettlementHistory(ctx, types.OptionSettlementHistoryQuery{}, types.PageRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "count:", len(settlements))
	fmt.Printf("%-30s %v\n", "page count:", page.Count)
}
