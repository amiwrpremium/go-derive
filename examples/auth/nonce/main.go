// Generates strictly-increasing nonces.
package main

import (
	"fmt"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	g := auth.NewNonceGen()
	for i := 0; i < 5; i++ {
		fmt.Printf("%-30s %v\n", "nonce:", g.Next())
	}
}
