// Computes the EIP-712 hashStruct of an ActionData.
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	a := auth.ActionData{
		SubaccountID: 1,
		Nonce:        42,
		Module:       common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Expiry:       1_700_000_000,
	}
	fmt.Printf("%-30s %v\n", "hash:", hex.EncodeToString(a.Hash()))
}
