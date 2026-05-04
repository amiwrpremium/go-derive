// Fetches details for one instrument by name.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	inst, err := c.GetInstrument(ctx, example.Instrument())
	example.Fatal(err)
	example.Print("name", inst.Name)
	example.Print("type", inst.Type)
	example.Print("tick size", inst.TickSize)
}
