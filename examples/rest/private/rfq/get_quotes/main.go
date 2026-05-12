// Paginates the configured subaccount's quotes — own quotes either
// active or in any historical state.
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

	quotes, page, err := c.GetQuotes(ctx, types.QuotesQuery{}, types.PageRequest{PageSize: 10})
	example.Fatal(err)
	example.Print("count", len(quotes))
	example.Print("total pages", page.NumPages)
}
