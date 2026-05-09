// Lists BTC perp instruments over WebSocket.
package main

import "github.com/amiwrpremium/go-derive"

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	insts, err := c.GetInstruments(ctx, "BTC", derive.InstrumentTypePerp)
	example.Fatal(err)
	example.Print("BTC perp count", len(insts))
}
