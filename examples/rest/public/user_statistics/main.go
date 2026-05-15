// One wallet's trading statistics. Required env: DERIVE_WALLET.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	wallet := os.Getenv("DERIVE_WALLET")
	if wallet == "" {
		log.Fatal("DERIVE_WALLET required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s, err := c.GetUserStatistics(ctx, types.UserStatisticsQuery{Wallet: wallet})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "total_fees:", s.TotalFees.String())
	fmt.Printf("%-30s %v\n", "total_trades:", s.TotalTrades)
	fmt.Printf("%-30s %v\n", "total_notional_volume:", s.TotalNotionalVolume.String())
	fmt.Printf("%-30s %v\n", "first_trade_ms:", s.FirstTradeTimestamp.Millis())
	fmt.Printf("%-30s %v\n", "last_trade_ms:", s.LastTradeTimestamp.Millis())
}
