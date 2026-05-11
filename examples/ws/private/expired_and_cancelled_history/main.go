// Exports the configured subaccount's expired and cancelled orders
// over WebSocket.
//
// Required env: DERIVE_WALLET.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	wallet := os.Getenv("DERIVE_WALLET")
	if wallet == "" {
		log.Fatal("DERIVE_WALLET required")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.ExpiredAndCancelledHistory(ctx, types.ExpiredAndCancelledHistoryInput{
		Wallet:    wallet,
		ExpirySec: 3600,
	})
	example.Fatal(err)
	example.Print("urls", len(res.PresignedURLs))
}
