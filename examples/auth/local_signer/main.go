// Builds an auth.LocalSigner from the configured private key.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	s := example.MustSigner()
	example.Print("address", s.Address())
	example.Print("owner", s.Owner())
}
