// Paginates the public trade tape — page-sized 5, prints the latest 5
// trades plus the total record count and page count Derive reports.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	trades, page, err := c.GetPublicTradeHistory(ctx, example.Instrument(),
		derive.PageRequest{PageSize: 5})
	example.Fatal(err)
	example.Print("trades returned", len(trades))
	example.Print("total records", page.Count)
	example.Print("total pages", page.NumPages)
	for _, t := range trades {
		example.Print(string(t.Direction)+" "+t.TradeAmount.String(), t.TradePrice)
	}
}
