// Lists active BTC perpetual instruments.
//
// `public/get_instruments` returns *static* instrument metadata
// (tick size, min/max amount, base/quote currencies, …) — it does not
// include live mark or index prices. Use `public/get_ticker` per
// instrument when you need those.
package main

import "github.com/amiwrpremium/go-derive"

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	insts, err := c.GetInstruments(ctx, "BTC", derive.InstrumentTypePerp)
	example.Fatal(err)
	example.Print("BTC perp count", len(insts))
	for i, in := range insts {
		if i >= 3 {
			break
		}
		example.Print(in.Name+" tick", in.TickSize)
		example.Print(in.Name+" min", in.MinimumAmount)
		example.Print(in.Name+" active", in.IsActive)
	}
}
