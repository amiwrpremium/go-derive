// Fetches the latest oracle signed-feed snapshot. Prints the per-currency
// summary (the full payload also carries oracle signatures, which this
// demo skips).
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
)

func main() {
	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork)
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	feeds, err := c.GetLatestSignedFeeds(ctx, "", 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-30s %v\n", "currencies with feeds:", len(feeds.SpotData))
	keys := make([]string, 0, len(feeds.SpotData))
	for k := range feeds.SpotData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "  "+k+" timestamp:", feeds.SpotData[k].Timestamp.Millis())
		fmt.Printf("%-30s %v\n", "  "+k+" price:", feeds.SpotData[k].Price.String())
	}
}
