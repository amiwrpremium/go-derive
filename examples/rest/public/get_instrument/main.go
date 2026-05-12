// Fetches details for one instrument by name.
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

	inst, err := c.GetInstrument(ctx, instrument)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "name:", inst.Name)
	fmt.Printf("%-30s %v\n", "type:", inst.Type)
	fmt.Printf("%-30s %v\n", "tick size:", inst.TickSize)
}
