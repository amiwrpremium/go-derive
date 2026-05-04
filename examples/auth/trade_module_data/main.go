// Hashes the per-trade module payload that goes into ActionData.Data.
package main

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	t := auth.TradeModuleData{
		Asset:       common.HexToAddress("0x1111111111111111111111111111111111111111"),
		LimitPrice:  decimal.RequireFromString("65000"),
		Amount:      decimal.RequireFromString("0.1"),
		MaxFee:      decimal.RequireFromString("10"),
		RecipientID: example.Subaccount(),
		IsBid:       true,
	}
	h, err := t.Hash()
	example.Fatal(err)
	example.Print("hash", hex.EncodeToString(h[:]))
}
