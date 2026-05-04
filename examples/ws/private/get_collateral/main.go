// Fetches collateral over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	cols, err := c.GetCollateral(ctx)
	example.Fatal(err)
	example.Print("count", len(cols))
}
