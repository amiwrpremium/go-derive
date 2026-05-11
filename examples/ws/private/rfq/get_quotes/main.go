// Paginates the configured subaccount's quotes over WebSocket.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	quotes, _, err := c.GetQuotes(ctx, types.QuotesQuery{}, types.PageRequest{PageSize: 10})
	example.Fatal(err)
	example.Print("count", len(quotes))
}
