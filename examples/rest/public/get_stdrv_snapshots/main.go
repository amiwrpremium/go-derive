// One wallet's staked-DRV balance snapshots over a time window.
//
// Required env: DERIVE_WALLET, DERIVE_FROM_SEC, DERIVE_TO_SEC.
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
	wallet := os.Getenv("DERIVE_WALLET")
	if wallet == "" {
		log.Fatal("DERIVE_WALLET required")
	}
	fromStr := os.Getenv("DERIVE_FROM_SEC")
	toStr := os.Getenv("DERIVE_TO_SEC")
	if fromStr == "" || toStr == "" {
		log.Fatal("DERIVE_FROM_SEC and DERIVE_TO_SEC required (Unix seconds)")
	}
	from, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_FROM_SEC=%q: %v", fromStr, err)
	}
	to, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_TO_SEC=%q: %v", toStr, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.GetStDRVSnapshots(ctx, types.STDRVSnapshotsQuery{
		Wallet:  wallet,
		FromSec: from,
		ToSec:   to,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "wallet:", res.Wallet)
	fmt.Printf("%-30s %v\n", "snapshots:", len(res.Snapshots))
	for i, s := range res.Snapshots {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "at sec:", s.TimestampSec)
		fmt.Printf("%-30s %v\n", "  amount:", s.Amount.String())
	}
}
