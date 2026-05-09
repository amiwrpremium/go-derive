// Paginates past orders over WebSocket.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	orders, _, err := c.GetOrderHistory(ctx, derive.PageRequest{PageSize: 10})
	example.Fatal(err)
	example.Print("count", len(orders))
}
