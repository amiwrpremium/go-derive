// Uses only the c.REST client from the facade.
package main

import "github.com/amiwrpremium/go-derive"

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustDerivePublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	insts, err := c.REST.GetInstruments(ctx, "BTC", derive.InstrumentTypePerp)
	example.Fatal(err)
	example.Print("BTC perps", len(insts))
}
