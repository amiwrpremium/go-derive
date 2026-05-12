// Simulates a margin calculation for a hypothetical portfolio over
// WebSocket. Public — no auth required.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wsNetwork := ws.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		wsNetwork = ws.WithMainnet()
	}
	c, err := ws.New(wsNetwork)
	if err != nil {
		log.Fatalf("ws.New: %v", err)
	}
	defer c.Close()
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("ws.Connect: %v", err)
	}
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
	fmt.Printf("%-30s %v\n", "post_initial_margin:", res.PostInitialMargin.String())
}
