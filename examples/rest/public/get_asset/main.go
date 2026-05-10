// Fetches one Asset record by name. Required env:
// DERIVE_ASSET_NAME.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	name := os.Getenv("DERIVE_ASSET_NAME")
	if name == "" {
		log.Fatal("DERIVE_ASSET_NAME required (e.g. USDC)")
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	a, err := c.GetAsset(ctx, name)
	example.Fatal(err)
	example.Print("asset_name", a.AssetName)
	example.Print("asset_type", string(a.AssetType))
	example.Print("currency", a.Currency)
	example.Print("is_collateral", a.IsCollateral)
	example.Print("is_position", a.IsPosition)
}
