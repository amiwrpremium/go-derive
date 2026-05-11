// Lists historical subaccount-value snapshots over WebSocket.
// Required env: DERIVE_FROM_MS, DERIVE_TO_MS, DERIVE_PERIOD_SEC.
package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	fromStr := os.Getenv("DERIVE_FROM_MS")
	toStr := os.Getenv("DERIVE_TO_MS")
	periodStr := os.Getenv("DERIVE_PERIOD_SEC")
	if fromStr == "" || toStr == "" || periodStr == "" {
		log.Fatal("DERIVE_FROM_MS, DERIVE_TO_MS and DERIVE_PERIOD_SEC required")
	}
	from, err := strconv.ParseInt(fromStr, 10, 64)
	example.Fatal(err)
	to, err := strconv.ParseInt(toStr, 10, 64)
	example.Fatal(err)
	period, err := strconv.ParseInt(periodStr, 10, 64)
	example.Fatal(err)

	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	subID, history, err := c.GetSubaccountValueHistory(ctx, types.SubaccountValueHistoryQuery{
		HistoryWindow: types.HistoryWindow{
			StartTimestamp: types.NewMillisTime(time.UnixMilli(from)),
			EndTimestamp:   types.NewMillisTime(time.UnixMilli(to)),
		},
		PeriodSec: period,
	})
	example.Fatal(err)
	example.Print("subaccount_id", subID)
	example.Print("snapshots", len(history))
}
