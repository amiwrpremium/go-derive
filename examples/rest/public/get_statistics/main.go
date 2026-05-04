// Fetches rolling daily / lifetime statistics (volume, fees, trades, OI)
// for one instrument and prints the headline numbers.
package main

import (
	"encoding/json"

	"github.com/amiwrpremium/go-derive/examples/example"
)

type statsResp struct {
	DailyNotionalVolume string `json:"daily_notional_volume"`
	DailyPremiumVolume  string `json:"daily_premium_volume"`
	DailyFees           string `json:"daily_fees"`
	DailyTrades         int    `json:"daily_trades"`
	TotalNotionalVolume string `json:"total_notional_volume"`
	TotalTrades         int    `json:"total_trades"`
	OpenInterest        string `json:"open_interest"`
}

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	raw, err := c.GetStatistics(ctx, example.Instrument())
	example.Fatal(err)

	var s statsResp
	example.Fatal(json.Unmarshal(raw, &s))
	example.Print("daily notional volume", s.DailyNotionalVolume)
	example.Print("daily trades", s.DailyTrades)
	example.Print("daily fees", s.DailyFees)
	example.Print("total trades", s.TotalTrades)
	example.Print("open interest", s.OpenInterest)
}
