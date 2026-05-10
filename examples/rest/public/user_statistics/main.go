// One wallet's trading statistics. Required env: DERIVE_WALLET.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	wallet := os.Getenv("DERIVE_WALLET")
	if wallet == "" {
		log.Fatal("DERIVE_WALLET required")
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	s, err := c.GetUserStatistics(ctx, map[string]any{"wallet": wallet})
	example.Fatal(err)
	example.Print("total_fees", s.TotalFees.String())
	example.Print("total_trades", s.TotalTrades)
	example.Print("total_notional_volume", s.TotalNotionalVolume.String())
	example.Print("first_trade_ms", s.FirstTradeTimestamp.Millis())
	example.Print("last_trade_ms", s.LastTradeTimestamp.Millis())
}
