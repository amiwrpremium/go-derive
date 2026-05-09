// Fetches the configured subaccount's past liquidation events.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	auctions, err := c.GetLiquidationHistory(ctx, nil)
	example.Fatal(err)
	example.Print("auction count", len(auctions))
	if len(auctions) > 0 {
		a := auctions[0]
		example.Print("first auction", a.AuctionID)
		example.Print("first auction type", a.AuctionType)
		example.Print("first auction fee", a.Fee.String())
	}
}
