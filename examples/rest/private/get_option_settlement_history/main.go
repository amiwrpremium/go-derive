// Lists option settlements for the configured subaccount.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	settlements, err := c.GetOptionSettlementHistory(ctx, nil)
	example.Fatal(err)
	example.Print("count", len(settlements))
}
