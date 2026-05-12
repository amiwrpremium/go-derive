// Per-wallet trading statistics for every wallet matching the
// supplied filters.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	stats, err := c.GetAllUserStatistics(ctx, 0)
	example.Fatal(err)
	example.Print("rows", len(stats))
	for i, s := range stats {
		if i >= 5 {
			break
		}
		example.Print("wallet", s.Wallet)
		example.Print("  total_fees", s.TotalFees.String())
		example.Print("  total_trades", s.TotalTrades)
	}
}
