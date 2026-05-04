// Fetches server time over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	t, err := c.GetTime(ctx)
	example.Fatal(err)
	example.Print("server time (ms)", t)
}
