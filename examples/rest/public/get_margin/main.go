// Simulates a margin calculation for a hypothetical portfolio.
// Public — no auth required. Pass the portfolio shape via
// types.PublicMarginInput per the docs at
// https://docs.derive.xyz/reference/public-get_margin.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetPublicMargin(ctx, types.PublicMarginInput{
		MarginType: enums.MarginTypePM,
		Market:     "BTC",
		SimulatedCollaterals: []types.SimulatedCollateral{
			{AssetName: "USDC", Amount: types.MustDecimal("10000")},
		},
		SimulatedPositions: []types.SimulatedPosition{},
	})
	example.Fatal(err)
	example.Print("pre_initial_margin", res.PreInitialMargin.String())
	example.Print("post_initial_margin", res.PostInitialMargin.String())
	example.Print("is_valid_trade", res.IsValidTrade)
}
