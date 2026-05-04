// Fetches a ticker (top-of-book + marks).
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	t, err := c.GetTicker(ctx, example.Instrument())
	example.Fatal(err)
	example.Print("instrument", t.InstrumentName)
	example.Print("bid", t.BestBidPrice)
	example.Print("ask", t.BestAskPrice)
	example.Print("mark", t.MarkPrice)
}
