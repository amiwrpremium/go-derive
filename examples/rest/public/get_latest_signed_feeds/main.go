// Fetches the latest oracle signed-feed snapshot. Prints the per-currency
// summary (the full payload also carries oracle signatures, which this
// demo skips).
package main

import (
	"sort"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	feeds, err := c.GetLatestSignedFeeds(ctx, "", 0)
	example.Fatal(err)

	example.Print("currencies with feeds", len(feeds.SpotData))
	keys := make([]string, 0, len(feeds.SpotData))
	for k := range feeds.SpotData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i >= 5 {
			break
		}
		example.Print("  "+k+" timestamp", feeds.SpotData[k].Timestamp.Millis())
		example.Print("  "+k+" price", feeds.SpotData[k].Price.String())
	}
}
