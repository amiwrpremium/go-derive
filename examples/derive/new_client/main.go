// Constructs the top-level facade with the configured network.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustDerivePublic()
	defer c.Close()
	example.Print("network", c.Network().Network)
	example.Print("chain id", c.Network().ChainID)
}
