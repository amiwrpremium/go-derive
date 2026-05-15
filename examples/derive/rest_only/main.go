// Uses only the c.REST client from the facade.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/derive"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	network := derive.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		network = derive.WithMainnet()
	}
	c, err := derive.NewClient(network)
	if err != nil {
		log.Fatalf("derive.NewClient: %v", err)
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	insts, err := c.REST.GetInstruments(ctx, types.InstrumentsQuery{Currency: "BTC", Kind: enums.InstrumentTypePerp})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "BTC perps:", len(insts))
}
