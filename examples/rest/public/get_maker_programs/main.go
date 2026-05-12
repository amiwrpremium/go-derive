// Lists active maker incentive programs. Each program has its own
// epoch, the asset types and currencies it covers, and the rewards
// paid out.
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

	programs, err := c.GetMakerPrograms(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "program count:", len(programs))
	for i, p := range programs {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "program:", p.Name)
		fmt.Printf("%-30s %v\n", "  asset_types:", p.AssetTypes)
		fmt.Printf("%-30s %v\n", "  currencies:", p.Currencies)
		fmt.Printf("%-30s %v\n", "  min_notional:", p.MinNotional.String())
	}
}
