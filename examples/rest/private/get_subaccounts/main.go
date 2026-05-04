// Lists every subaccount the wallet owns.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	ids, err := c.GetSubaccounts(ctx)
	example.Fatal(err)
	example.Print("subaccount ids", ids)
}
