// Lists Asset records matching the supplied filter.
//
// Required env: DERIVE_CURRENCY (the underlying, e.g. "ETH").
// Optional env: DERIVE_ASSET_TYPE (defaults to "erc20"; one of
// "erc20", "option", "perp") and DERIVE_INCLUDE_EXPIRED=1.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func main() {
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		log.Fatal("DERIVE_CURRENCY required")
	}
	assetType := enums.AssetType(os.Getenv("DERIVE_ASSET_TYPE"))
	if assetType == "" {
		assetType = enums.AssetTypeERC20
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	assets, err := c.GetAssets(ctx, assetType, currency, os.Getenv("DERIVE_INCLUDE_EXPIRED") == "1")
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
