// Generates strictly-increasing nonces.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	g := derive.NewNonceGen()
	for i := 0; i < 5; i++ {
		example.Print("nonce", g.Next())
	}
}
