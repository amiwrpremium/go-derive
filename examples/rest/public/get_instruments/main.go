// Lists active BTC perpetual instruments.
//
// `public/get_instruments` returns *static* instrument metadata
// (tick size, min/max amount, base/quote currencies, …) — it does not
// include live mark or index prices. Use `public/get_ticker` per
// instrument when you need those.
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

	insts, err := c.GetInstruments(ctx, types.InstrumentsQuery{Currency: "BTC", Kind: enums.InstrumentTypePerp})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "BTC perp count:", len(insts))
	for i, in := range insts {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", in.Name+" tick:", in.TickSize)
		fmt.Printf("%-30s %v\n", in.Name+" min:", in.MinimumAmount)
		fmt.Printf("%-30s %v\n", in.Name+" active:", in.IsActive)
	}
}
