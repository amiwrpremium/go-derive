// Calculates the mark-to-market and maintenance-margin snapshot for
// one subaccount. Required env: DERIVE_SUBACCOUNT (any subaccount
// id, not necessarily the configured signer's).
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	idStr := os.Getenv("DERIVE_SUBACCOUNT")
	if idStr == "" {
		log.Fatal("DERIVE_SUBACCOUNT required")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_SUBACCOUNT=%q: %v", idStr, err)
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	snap, err := c.MarginWatch(ctx, id, false, false)
	example.Fatal(err)
	example.Print("subaccount_id", snap.SubaccountID)
	example.Print("margin_type", snap.MarginType)
	example.Print("subaccount_value", snap.SubaccountValue.String())
	example.Print("maintenance_margin", snap.MaintenanceMargin.String())
	example.Print("collaterals", len(snap.Collaterals))
	example.Print("positions", len(snap.Positions))
}
