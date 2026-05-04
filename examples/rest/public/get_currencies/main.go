// Lists all quote currencies supported on the configured network.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	got, err := c.GetCurrencies(ctx)
	example.Fatal(err)
	example.Print("currencies", got)
}
