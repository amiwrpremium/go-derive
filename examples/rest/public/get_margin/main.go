// Simulates a margin calculation for a hypothetical portfolio.
// Public — no auth required. Pass the portfolio shape via
// types.PublicMarginInput per the docs at
// https://docs.derive.xyz/reference/public-get_margin.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/types"
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.GetPublicMargin(ctx, types.PublicMarginInput{
		MarginType: enums.MarginTypePM,
		Market:     "BTC",
		SimulatedCollaterals: []types.SimulatedCollateral{
			{AssetName: "USDC", Amount: types.MustDecimal("10000")},
		},
		SimulatedPositions: []types.SimulatedPosition{},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "pre_initial_margin:", res.PreInitialMargin.String())
	fmt.Printf("%-30s %v\n", "post_initial_margin:", res.PostInitialMargin.String())
	fmt.Printf("%-30s %v\n", "is_valid_trade:", res.IsValidTrade)
}
