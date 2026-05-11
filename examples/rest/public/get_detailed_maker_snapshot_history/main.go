// Per-quote maker-snapshot rows for one program / epoch.
//
// Required env: DERIVE_PROGRAM_NAME, DERIVE_EPOCH_START (Unix
// milliseconds), DERIVE_WALLET (the maker to scope to).
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	name := os.Getenv("DERIVE_PROGRAM_NAME")
	if name == "" {
		log.Fatal("DERIVE_PROGRAM_NAME required")
	}
	epochStr := os.Getenv("DERIVE_EPOCH_START")
	if epochStr == "" {
		log.Fatal("DERIVE_EPOCH_START required (Unix ms)")
	}
	epoch, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_EPOCH_START=%q: %v", epochStr, err)
	}
	wallet := os.Getenv("DERIVE_WALLET")
	if wallet == "" {
		log.Fatal("DERIVE_WALLET required")
	}

	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetDetailedMakerSnapshotHistory(ctx, types.DetailedMakerSnapshotHistoryQuery{
		ProgramName:         name,
		EpochStartTimestamp: epoch,
		Wallet:              wallet,
	}, types.PageRequest{})
	example.Fatal(err)
	example.Print("program", res.Program.Name)
	example.Print("snapshots", len(res.Snapshots))
	example.Print("total pages", res.Pagination.NumPages)
	for i, s := range res.Snapshots {
		if i >= 3 {
			break
		}
		example.Print("snapshot", s.InstrumentName)
		example.Print("  notional", s.Notional.String())
		example.Print("  scaled_notional", s.ScaledNotional.String())
	}
}
