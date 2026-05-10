// Fetches historical USDC borrow / supply APY prints over the last
// 24 hours.
package main

import (
	"time"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	end := time.Now().Unix()
	start := end - 24*3600

	rates, page, err := c.GetInterestRateHistory(ctx, map[string]any{
		"from_timestamp_sec": start,
		"to_timestamp_sec":   end,
	})
	example.Fatal(err)
	example.Print("rate prints", len(rates))
	example.Print("total pages", page.NumPages)
	for i, r := range rates {
		if i >= 3 {
			break
		}
		example.Print("at sec", r.TimestampSec)
		example.Print("  borrow / supply APY", r.BorrowAPY.String()+" / "+r.SupplyAPY.String())
	}
}
