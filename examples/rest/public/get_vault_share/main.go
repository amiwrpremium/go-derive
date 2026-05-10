// Fetches per-block snapshots of one vault token's price-per-share
// over the last 24 hours.
//
// Required: DERIVE_VAULT_NAME (run get_vault_statistics first to
// discover available vault names).
package main

import (
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	name := os.Getenv("DERIVE_VAULT_NAME")
	if name == "" {
		log.Fatal("DERIVE_VAULT_NAME required (run get_vault_statistics to discover names)")
	}

	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	end := time.Now().Unix()
	start := end - 24*3600

	shares, page, err := c.GetVaultShare(ctx, map[string]any{
		"vault_name":         name,
		"from_timestamp_sec": start,
		"to_timestamp_sec":   end,
	})
	example.Fatal(err)
	example.Print("snapshot count", len(shares))
	example.Print("total pages", page.NumPages)
	for i, s := range shares {
		if i >= 3 {
			break
		}
		example.Print("block", s.BlockNumber)
		example.Print("  usd_value", s.USDValue.String())
		example.Print("  base_value", s.BaseValue.String())
	}
}
