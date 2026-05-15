// Lists historical subaccount-value snapshots. Required env:
// DERIVE_FROM_MS, DERIVE_TO_MS, DERIVE_PERIOD_SEC (one of 900,
// 3600, 86400, 604800).
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	subStr := os.Getenv("DERIVE_SUBACCOUNT")
	if subStr == "" {
		log.Fatal("DERIVE_SUBACCOUNT required")
	}
	subaccount, err := strconv.ParseInt(subStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_SUBACCOUNT=%q: %v", subStr, err)
	}
	key := os.Getenv("DERIVE_SESSION_KEY")
	if key == "" {
		log.Fatal("DERIVE_SESSION_KEY required")
	}
	var signer auth.Signer
	if owner := os.Getenv("DERIVE_OWNER"); owner != "" {
		signer, err = auth.NewSessionKeySigner(key, common.HexToAddress(owner))
	} else {
		signer, err = auth.NewLocalSigner(key)
	}
	if err != nil {
		log.Fatalf("signer: %v", err)
	}

	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork, rest.WithSigner(signer), rest.WithSubaccount(subaccount))
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	fromStr := os.Getenv("DERIVE_FROM_MS")
	toStr := os.Getenv("DERIVE_TO_MS")
	periodStr := os.Getenv("DERIVE_PERIOD_SEC")
	if fromStr == "" || toStr == "" || periodStr == "" {
		log.Fatal("DERIVE_FROM_MS, DERIVE_TO_MS and DERIVE_PERIOD_SEC required")
	}
	from, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	to, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	period, err := strconv.ParseInt(periodStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	subID, history, err := c.GetSubaccountValueHistory(ctx, types.SubaccountValueHistoryQuery{
		HistoryWindow: types.HistoryWindow{
			StartTimestamp: types.MillisTimeFromMillis(from),
			EndTimestamp:   types.MillisTimeFromMillis(to),
		},
		PeriodSec: period,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "subaccount_id:", subID)
	fmt.Printf("%-30s %v\n", "snapshots:", len(history))
}
