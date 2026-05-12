// Lists historical oracle spot-feed values over WebSocket.
// Required env: DERIVE_CURRENCY.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		log.Fatal("DERIVE_CURRENCY required")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	cur, items, err := c.GetSpotFeedHistory(ctx, types.SpotFeedHistoryQuery{
		Currency:  currency,
		PeriodSec: 3600,
	})
	example.Fatal(err)
	example.Print("currency", cur)
	example.Print("items", len(items))
}
