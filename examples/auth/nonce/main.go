// Generates strictly-increasing nonces.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	g := auth.NewNonceGen()
	for i := 0; i < 5; i++ {
		example.Print("nonce", g.Next())
	}
}
