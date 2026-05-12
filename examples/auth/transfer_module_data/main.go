// Hashes the per-transfer module payload.
package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	t := auth.TransferModuleData{
		ToSubaccount: 99,
		Asset:        common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Amount:       decimal.RequireFromString("10"),
	}
	h, err := t.Hash()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "hash:", hex.EncodeToString(h[:]))
}
