// Demonstrates Signature.Hex serialisation.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	var s auth.Signature
	for i := range s {
		s[i] = byte(i)
	}
	example.Print("hex", s.Hex())
}
