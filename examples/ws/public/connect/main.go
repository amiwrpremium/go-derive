// Opens a WebSocket connection and reports its state.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()
	example.Print("connected", c.IsConnected())
}
