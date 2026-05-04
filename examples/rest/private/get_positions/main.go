// Lists the subaccount's open positions.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	ps, err := c.GetPositions(ctx)
	example.Fatal(err)
	example.Print("count", len(ps))
	for _, p := range ps {
		example.Print(p.InstrumentName, p.Amount)
	}
}
