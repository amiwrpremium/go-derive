// Paginates past orders.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	orders, page, err := c.GetOrderHistory(ctx, types.PageRequest{PageSize: 10})
	example.Fatal(err)
	example.Print("count", len(orders))
	example.Print("total pages", page.NumPages)
}
