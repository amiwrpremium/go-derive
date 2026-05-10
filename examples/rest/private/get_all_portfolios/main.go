// Lists per-subaccount portfolio snapshots for every subaccount the
// configured signer's wallet owns. Each snapshot carries the full
// margin breakdown plus collateral / position / open-order arrays.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	portfolios, err := c.GetAllPortfolios(ctx)
	example.Fatal(err)
	example.Print("portfolio count", len(portfolios))
	for i, p := range portfolios {
		if i >= 3 {
			break
		}
		example.Print("subaccount", p.SubaccountID)
		example.Print("  margin_type", p.MarginType)
		example.Print("  subaccount_value", p.SubaccountValue.String())
		example.Print("  positions", len(p.Positions))
		example.Print("  open_orders", len(p.OpenOrders))
	}
}
