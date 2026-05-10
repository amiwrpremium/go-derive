// One page (up to 500 entries) of the points leaderboard for one
// program. Required env: DERIVE_PROGRAM_NAME.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	name := os.Getenv("DERIVE_PROGRAM_NAME")
	if name == "" {
		log.Fatal("DERIVE_PROGRAM_NAME required")
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetPointsLeaderboard(ctx, name, 1)
	example.Fatal(err)
	example.Print("pages", res.Pages)
	example.Print("total_users", res.TotalUsers)
	example.Print("entries on page", len(res.Leaderboard))
	for i, e := range res.Leaderboard {
		if i >= 5 {
			break
		}
		example.Print("rank", e.Rank)
		example.Print("  wallet", e.Wallet)
		example.Print("  points", e.Points.String())
	}
}
