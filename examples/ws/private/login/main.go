// Authenticates a WebSocket session via public/login.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx) // helper does Connect + Login.
	defer c.Close()
	example.Print("logged in", c.IsConnected())
}
