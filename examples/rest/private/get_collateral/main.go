// Lists the subaccount's collateral positions.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	cols, err := c.GetCollateral(ctx)
	example.Fatal(err)
	for _, c := range cols {
		example.Print(c.AssetName, c.Amount)
	}
}
