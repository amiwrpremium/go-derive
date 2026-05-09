// Paginates the public trade tape over WebSocket.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	trades, _, err := c.GetPublicTradeHistory(ctx, example.Instrument(),
		derive.PageRequest{PageSize: 5})
	example.Fatal(err)
	example.Print("count", len(trades))
}
