// Paginates the user's filled trades.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	trades, page, err := c.GetTradeHistory(ctx, derive.PageRequest{PageSize: 10})
	example.Fatal(err)
	example.Print("count", len(trades))
	example.Print("total pages", page.NumPages)
}
