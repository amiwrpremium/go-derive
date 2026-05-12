// Lists the network-wide liquidation auction history (paginated).
// Counterpart to private/get_liquidation_history (single-subaccount).
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	auctions, page, err := c.GetPublicLiquidationHistory(ctx, types.LiquidationHistoryQuery{}, types.PageRequest{})
	example.Fatal(err)
	example.Print("auctions", len(auctions))
	example.Print("total pages", page.NumPages)
	for i, a := range auctions {
		if i >= 3 {
			break
		}
		example.Print("auction", a.AuctionID)
		example.Print("  type", string(a.AuctionType))
		example.Print("  subaccount", a.SubaccountID)
	}
}
