// One wallet's staked-DRV balance snapshots over a time window.
//
// Required env: DERIVE_WALLET, DERIVE_FROM_SEC, DERIVE_TO_SEC.
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
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
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetStDRVSnapshots(ctx, map[string]any{
		"wallet":   wallet,
		"from_sec": from,
		"to_sec":   to,
	})
	example.Fatal(err)
	example.Print("wallet", res.Wallet)
	example.Print("snapshots", len(res.Snapshots))
	for i, s := range res.Snapshots {
		if i >= 3 {
			break
		}
		example.Print("at sec", s.TimestampSec)
		example.Print("  amount", s.Amount.String())
	}
}
