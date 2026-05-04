// Lists currently-open orders.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	orders, err := c.GetOpenOrders(ctx)
	example.Fatal(err)
	example.Print("count", len(orders))
}
