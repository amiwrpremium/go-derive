// Cancels every open order on the subaccount over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	n, err := c.CancelAll(ctx)
	example.Fatal(err)
	example.Print("cancelled", n)
}
