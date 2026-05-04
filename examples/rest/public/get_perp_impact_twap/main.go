// Fetches the time-weighted average impact price for one currency's
// perpetual book over the last hour.
//
// Derive returns a small object with the mid-price diff, ask-impact diff,
// and bid-impact diff (all decimal strings). When the book is quiet, all
// three are "0".
package main

import (
	"encoding/json"
	"time"

	"github.com/amiwrpremium/go-derive/examples/example"
)

type impactResp struct {
	Currency        string `json:"currency"`
	MidDiff         string `json:"mid_price_diff_twap"`
	AskImpactDiff   string `json:"ask_impact_diff_twap"`
	BidImpactDiff   string `json:"bid_impact_diff_twap"`
}

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	now := time.Now()
	raw, err := c.GetPerpImpactTWAP(ctx, map[string]any{
		"currency":   "BTC",
		"start_time": now.Add(-time.Hour).UnixMilli(),
		"end_time":   now.UnixMilli(),
	})
	example.Fatal(err)

	var r impactResp
	example.Fatal(json.Unmarshal(raw, &r))
	example.Print("currency", r.Currency)
	example.Print("mid diff TWAP", r.MidDiff)
	example.Print("ask impact TWAP", r.AskImpactDiff)
	example.Print("bid impact TWAP", r.BidImpactDiff)
}
