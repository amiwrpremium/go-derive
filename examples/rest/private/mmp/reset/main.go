// Unfreezes market-maker protection for one currency.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	example.Fatal(c.ResetMMP(ctx, "BTC"))
	example.Print("reset", "ok")
}
