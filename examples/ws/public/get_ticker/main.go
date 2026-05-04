// Fetches a ticker over WebSocket (single RPC, not a subscription).
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	t, err := c.GetTicker(ctx, example.Instrument())
	example.Fatal(err)
	example.Print("mark", t.MarkPrice)
}
