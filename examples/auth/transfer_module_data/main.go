// Hashes the per-transfer module payload.
package main

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	t := derive.TransferModuleData{
		ToSubaccount: 99,
		Asset:        common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Amount:       decimal.RequireFromString("10"),
	}
	h, err := t.Hash()
	example.Fatal(err)
	example.Print("hash", hex.EncodeToString(h[:]))
}
