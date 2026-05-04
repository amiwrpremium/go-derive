// Fetches historical funding-rate prints for one perpetual instrument
// and prints the most recent few entries.
package main

import (
	"encoding/json"

	"github.com/amiwrpremium/go-derive/examples/example"
)

type fundingResp struct {
	History []struct {
		Timestamp   int64  `json:"timestamp"`
		FundingRate string `json:"funding_rate"`
	} `json:"funding_rate_history"`
}

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	raw, err := c.GetFundingRateHistory(ctx, map[string]any{
		"instrument_name": example.Instrument(),
	})
	example.Fatal(err)

	var r fundingResp
	example.Fatal(json.Unmarshal(raw, &r))
	example.Print("entries returned", len(r.History))
	for i, e := range r.History {
		if i >= 5 {
			break
		}
		example.Print("rate at ms", e.Timestamp)
		example.Print("  funding_rate", e.FundingRate)
	}
}
