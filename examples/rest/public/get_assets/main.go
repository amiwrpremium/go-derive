// Lists Asset records matching the supplied filter.
//
// Required env: DERIVE_CURRENCY (the underlying, e.g. "ETH").
// Optional env: DERIVE_ASSET_TYPE (defaults to "erc20"; one of
// "erc20", "option", "perp") and DERIVE_INCLUDE_EXPIRED=1.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/rest"
)

func main() {
	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork)
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		log.Fatal("DERIVE_CURRENCY required")
	}
	assetType := enums.AssetType(os.Getenv("DERIVE_ASSET_TYPE"))
	if assetType == "" {
		assetType = enums.AssetTypeERC20
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	assets, err := c.GetAssets(ctx, assetType, currency, os.Getenv("DERIVE_INCLUDE_EXPIRED") == "1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "assets:", len(assets))
	for i, a := range assets {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "asset:", a.AssetName)
		fmt.Printf("%-30s %v\n", "  type:", string(a.AssetType))
	}
}
