// Resets MMP for one currency over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	example.Fatal(c.ResetMMP(ctx, "BTC"))
	example.Print("reset", "ok")
}
