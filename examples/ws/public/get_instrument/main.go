// Fetches one instrument over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	inst, err := c.GetInstrument(ctx, example.Instrument())
	example.Fatal(err)
	example.Print("name", inst.Name)
}
