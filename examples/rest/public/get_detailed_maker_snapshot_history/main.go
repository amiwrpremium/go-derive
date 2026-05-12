// Per-quote maker-snapshot rows for one program / epoch.
//
// Required env: DERIVE_PROGRAM_NAME, DERIVE_EPOCH_START (Unix
// milliseconds), DERIVE_WALLET (the maker to scope to).
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.GetDetailedMakerSnapshotHistory(ctx, types.DetailedMakerSnapshotHistoryQuery{
		ProgramName:         name,
		EpochStartTimestamp: epoch,
		Wallet:              wallet,
	}, types.PageRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "program:", res.Program.Name)
	fmt.Printf("%-30s %v\n", "snapshots:", len(res.Snapshots))
	fmt.Printf("%-30s %v\n", "total pages:", res.Pagination.NumPages)
	for i, s := range res.Snapshots {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "snapshot:", s.InstrumentName)
		fmt.Printf("%-30s %v\n", "  notional:", s.Notional.String())
		fmt.Printf("%-30s %v\n", "  scaled_notional:", s.ScaledNotional.String())
	}
}
