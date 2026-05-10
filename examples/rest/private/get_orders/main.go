// Paginates orders on the configured subaccount, optionally filtered
// by instrument / label / status. Pass nil filter to omit filters
// and page through every order on the subaccount.
//
// Counterpart to GetOrderHistory (time-window-based pagination).
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

	orders, page, err := c.GetOrders(ctx, types.PageRequest{PageSize: 10}, nil)
	example.Fatal(err)
	example.Print("orders", len(orders))
	example.Print("total pages", page.NumPages)
	for i, o := range orders {
		if i >= 3 {
			break
		}
		example.Print("order", o.OrderID)
		example.Print("  status", string(o.OrderStatus))
		example.Print("  instrument", o.InstrumentName)
	}
}
