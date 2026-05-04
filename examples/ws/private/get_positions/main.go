// Fetches positions over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	ps, err := c.GetPositions(ctx)
	example.Fatal(err)
	example.Print("count", len(ps))
}
