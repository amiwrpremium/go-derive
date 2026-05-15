// Calculates the mark-to-market and maintenance-margin snapshot for
// one subaccount. Required env: DERIVE_SUBACCOUNT (any subaccount
// id, not necessarily the configured signer's).
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
	idStr := os.Getenv("DERIVE_SUBACCOUNT")
	if idStr == "" {
		log.Fatal("DERIVE_SUBACCOUNT required")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_SUBACCOUNT=%q: %v", idStr, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	snap, err := c.MarginWatch(ctx, types.MarginWatchQuery{SubaccountID: id})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "subaccount_id:", snap.SubaccountID)
	fmt.Printf("%-30s %v\n", "margin_type:", snap.MarginType)
	fmt.Printf("%-30s %v\n", "subaccount_value:", snap.SubaccountValue.String())
	fmt.Printf("%-30s %v\n", "maintenance_margin:", snap.MaintenanceMargin.String())
	fmt.Printf("%-30s %v\n", "collaterals:", len(snap.Collaterals))
	fmt.Printf("%-30s %v\n", "positions:", len(snap.Positions))
}
