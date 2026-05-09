// Demonstrates Signature.Hex serialisation.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	var s derive.Signature
	for i := range s {
		s[i] = byte(i)
	}
	example.Print("hex", s.Hex())
}
