// Lists every active perpetual instrument across all currencies,
// paginated. Counterpart to GetInstruments which filters by
// currency.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/enums"
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

	insts, page, err := c.GetAllInstruments(ctx, types.AllInstrumentsQuery{Kind: enums.InstrumentTypePerp}, types.PageRequest{PageSize: 50})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "instrument count:", len(insts))
	fmt.Printf("%-30s %v\n", "total pages:", page.NumPages)
	for i, in := range insts {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "instrument:", in.Name)
	}
}
