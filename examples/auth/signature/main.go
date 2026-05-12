// Demonstrates Signature.Hex serialisation.
package main

import (
	"fmt"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	var s auth.Signature
	for i := range s {
		s[i] = byte(i)
	}
	fmt.Printf("%-30s %v\n", "hex:", s.Hex())
}
