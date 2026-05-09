// Paginates deposits over WebSocket.
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

	events, _, err := c.GetDepositHistory(ctx, derive.PageRequest{PageSize: 10})
	example.Fatal(err)
	example.Print("count", len(events))
}
