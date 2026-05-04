// Cancels every open order on the subaccount.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	n, err := c.CancelAll(ctx)
	example.Fatal(err)
	example.Print("cancelled", n)
}
