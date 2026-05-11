// Lists option settlements for the configured subaccount over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	settlements, err := c.GetOptionSettlementHistory(ctx, nil)
	example.Fatal(err)
	example.Print("count", len(settlements))
}
