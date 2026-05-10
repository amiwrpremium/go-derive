// Lists auction bids placed by the configured subaccount as a
// liquidator (paginated).
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	bids, page, err := c.GetLiquidatorHistory(ctx, nil)
	example.Fatal(err)
	example.Print("bids", len(bids))
	example.Print("total pages", page.NumPages)
	for i, b := range bids {
		if i >= 3 {
			break
		}
		example.Print("bid at ms", b.Timestamp.Millis())
		example.Print("  cash_received", b.CashReceived.String())
		example.Print("  realized_pnl", b.RealizedPnL.String())
	}
}
