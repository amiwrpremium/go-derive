// Hashes the per-trade module payload that goes into ActionData.Data.
package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/pkg/auth"
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

	t := auth.TradeModuleData{
		Asset:       common.HexToAddress("0x1111111111111111111111111111111111111111"),
		LimitPrice:  decimal.RequireFromString("65000"),
		Amount:      decimal.RequireFromString("0.1"),
		MaxFee:      decimal.RequireFromString("10"),
		RecipientID: subaccount,
		IsBid:       true,
	}
	h, err := t.Hash()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "hash:", hex.EncodeToString(h[:]))
}
