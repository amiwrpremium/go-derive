// Lists historical subaccount-value snapshots. Required env:
// DERIVE_FROM_SEC, DERIVE_TO_SEC, DERIVE_PERIOD (e.g. "1h").
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	fromStr := os.Getenv("DERIVE_FROM_SEC")
	toStr := os.Getenv("DERIVE_TO_SEC")
	period := os.Getenv("DERIVE_PERIOD")
	if fromStr == "" || toStr == "" || period == "" {
		log.Fatal("DERIVE_FROM_SEC, DERIVE_TO_SEC and DERIVE_PERIOD required")
	}
	from, err := strconv.ParseInt(fromStr, 10, 64)
	example.Fatal(err)
	to, err := strconv.ParseInt(toStr, 10, 64)
	example.Fatal(err)

	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	subID, history, err := c.GetSubaccountValueHistory(ctx, map[string]any{
		"from_timestamp_sec": from,
		"to_timestamp_sec":   to,
		"period":             period,
	})
	example.Fatal(err)
	example.Print("subaccount_id", subID)
	example.Print("snapshots", len(history))
}
