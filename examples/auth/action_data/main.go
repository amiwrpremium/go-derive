// Computes the EIP-712 hashStruct of an ActionData.
package main

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	a := auth.ActionData{
		SubaccountID: 1,
		Nonce:        42,
		Module:       common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Expiry:       1_700_000_000,
	}
	example.Print("hash", hex.EncodeToString(a.Hash()))
}
