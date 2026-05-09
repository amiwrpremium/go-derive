// Fetches rolling daily / lifetime statistics (volume, fees, trades, OI)
// for one instrument and prints the headline numbers.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	s, err := c.GetStatistics(ctx, example.Instrument())
	example.Fatal(err)

	example.Print("daily notional volume", s.DailyNotionalVolume.String())
	example.Print("daily trades", s.DailyTrades)
	example.Print("daily fees", s.DailyFees.String())
	example.Print("total trades", s.TotalTrades)
	example.Print("open interest", s.OpenInterest.String())
}
