// Fetches the per-wallet score breakdown for one maker incentive
// program at one epoch.
//
// Required: DERIVE_PROGRAM_NAME and DERIVE_EPOCH_START (Unix
// seconds). Run `go run ./examples/rest/public/get_maker_programs`
// first to discover available program names.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/types"
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.GetMakerProgramScores(ctx, types.MakerProgramScoresQuery{ProgramName: name, EpochStartTimestamp: epoch})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "program:", res.Program.Name)
	fmt.Printf("%-30s %v\n", "total_score:", res.TotalScore.String())
	fmt.Printf("%-30s %v\n", "total_volume:", res.TotalVolume.String())
	fmt.Printf("%-30s %v\n", "score count:", len(res.Scores))
	for i, s := range res.Scores {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "wallet:", s.Wallet.String())
		fmt.Printf("%-30s %v\n", "  total_score:", s.TotalScore.String())
	}
}
