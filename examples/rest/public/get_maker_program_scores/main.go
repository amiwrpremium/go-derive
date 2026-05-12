// Fetches the per-wallet score breakdown for one maker incentive
// program at one epoch.
//
// Required: DERIVE_PROGRAM_NAME and DERIVE_EPOCH_START (Unix
// seconds). Run `go run ./examples/rest/public/get_maker_programs`
// first to discover available program names.
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	name := os.Getenv("DERIVE_PROGRAM_NAME")
	if name == "" {
		log.Fatal("DERIVE_PROGRAM_NAME required (run get_maker_programs to discover names)")
	}
	epochStr := os.Getenv("DERIVE_EPOCH_START")
	if epochStr == "" {
		log.Fatal("DERIVE_EPOCH_START required (Unix seconds)")
	}
	epoch, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_EPOCH_START=%q: %v", epochStr, err)
	}

	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetMakerProgramScores(ctx, name, epoch)
	example.Fatal(err)
	example.Print("program", res.Program.Name)
	example.Print("total_score", res.TotalScore.String())
	example.Print("total_volume", res.TotalVolume.String())
	example.Print("score count", len(res.Scores))
	for i, s := range res.Scores {
		if i >= 3 {
			break
		}
		example.Print("wallet", s.Wallet.String())
		example.Print("  total_score", s.TotalScore.String())
	}
}
