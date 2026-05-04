// Uses only the c.WS client from the facade — connects and fetches a ticker.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustDerivePublic()
	defer c.Close()

	ctx, cancel := example.Timeout()
	defer cancel()
	example.Fatal(c.WS.Connect(ctx))

	tk, err := c.WS.GetTicker(ctx, example.Instrument())
	example.Fatal(err)
	example.Print("mark", tk.MarkPrice)
}
