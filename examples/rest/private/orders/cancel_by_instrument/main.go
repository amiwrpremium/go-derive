// Cancels every open order on one instrument.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	n, err := c.CancelByInstrument(ctx, example.Instrument())
	example.Fatal(err)
	example.Print("cancelled", n)
}
