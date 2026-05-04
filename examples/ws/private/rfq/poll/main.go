// Polls outstanding RFQs over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	rfqs, err := c.PollRFQs(ctx)
	example.Fatal(err)
	example.Print("count", len(rfqs))
}
