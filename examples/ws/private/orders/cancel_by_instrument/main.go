// Cancels orders on one instrument over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	n, err := c.CancelByInstrument(ctx, example.Instrument())
	example.Fatal(err)
	example.Print("cancelled", n)
}
