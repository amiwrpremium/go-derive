// Lists BTC perp instruments over WebSocket.
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
	insts, err := c.GetInstruments(ctx, types.InstrumentsQuery{Currency: "BTC", Kind: enums.InstrumentTypePerp})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "BTC perp count:", len(insts))
}
