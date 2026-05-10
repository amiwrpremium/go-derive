// Aggregate (currency, instrument_type) statistics across every
// instrument.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	stats, err := c.GetAllStatistics(ctx, 0)
	example.Fatal(err)
	example.Print("rows", len(stats))
	for i, s := range stats {
		if i >= 5 {
			break
		}
		example.Print("tuple", s.Currency+"/"+s.InstrumentType)
		example.Print("  daily_notional_volume", s.DailyNotionalVolume.String())
		example.Print("  open_interest", s.OpenInterest.String())
	}
}
