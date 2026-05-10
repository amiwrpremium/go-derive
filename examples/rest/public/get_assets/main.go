// Lists Asset records matching the supplied filter. Optional env:
// DERIVE_CURRENCY.
package main

import (
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	params := map[string]any{}
	if cur := os.Getenv("DERIVE_CURRENCY"); cur != "" {
		params["currency"] = cur
	}
	assets, err := c.GetAssets(ctx, params)
	example.Fatal(err)
	example.Print("assets", len(assets))
	for i, a := range assets {
		if i >= 5 {
			break
		}
		example.Print("asset", a.AssetName)
		example.Print("  type", string(a.AssetType))
	}
}
