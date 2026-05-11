// Simulates a margin calculation for a hypothetical portfolio over
// WebSocket. Public — no auth required.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	res, err := c.GetPublicMargin(ctx, types.PublicMarginInput{
		MarginType: enums.MarginTypePM,
		Market:     "BTC",
		SimulatedCollaterals: []types.SimulatedCollateral{
			{AssetName: "USDC", Amount: types.MustDecimal("10000")},
		},
		SimulatedPositions: []types.SimulatedPosition{},
	})
	example.Fatal(err)
	example.Print("post_initial_margin", res.PostInitialMargin.String())
}
