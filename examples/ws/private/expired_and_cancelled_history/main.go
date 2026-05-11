// Exports the configured subaccount's expired and cancelled orders
// over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.ExpiredAndCancelledHistory(ctx, nil)
	example.Fatal(err)
	example.Print("urls", len(res.PresignedURLs))
}
