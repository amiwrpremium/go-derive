// Fetches a single subaccount snapshot.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	sa, err := c.GetSubaccount(ctx)
	example.Fatal(err)
	example.Print("subaccount id", sa.SubaccountID)
	example.Print("equity", sa.SubaccountValue)
	example.Print("init margin", sa.InitialMargin)
}
