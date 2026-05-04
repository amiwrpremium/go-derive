// Fetches the latest oracle signed-feed snapshot. Prints the per-currency
// summary (the full payload also carries oracle signatures, which this
// demo skips).
package main

import (
	"encoding/json"
	"sort"

	"github.com/amiwrpremium/go-derive/examples/example"
)

type spotEntry struct {
	Currency  string `json:"currency"`
	Timestamp int64  `json:"timestamp"`
}

type signedFeeds struct {
	SpotData map[string]spotEntry `json:"spot_data"`
}

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	raw, err := c.GetLatestSignedFeeds(ctx, nil)
	example.Fatal(err)

	var sf signedFeeds
	example.Fatal(json.Unmarshal(raw, &sf))
	example.Print("currencies with feeds", len(sf.SpotData))
	keys := make([]string, 0, len(sf.SpotData))
	for k := range sf.SpotData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i >= 5 {
			break
		}
		example.Print("  "+k+" timestamp", sf.SpotData[k].Timestamp)
	}
}
