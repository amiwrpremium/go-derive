// Lists currencies over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	got, err := c.GetCurrencies(ctx)
	example.Fatal(err)
	example.Print("currencies", got)
}
