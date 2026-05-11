// Cancels every order on the configured subaccount that carries the
// given nonce, over WebSocket. Required env: DERIVE_INSTRUMENT,
// DERIVE_NONCE.
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	instrument := os.Getenv("DERIVE_INSTRUMENT")
	nonceStr := os.Getenv("DERIVE_NONCE")
	if instrument == "" || nonceStr == "" {
		log.Fatal("DERIVE_INSTRUMENT and DERIVE_NONCE required")
	}
	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	example.Fatal(err)

	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.CancelByNonce(ctx, instrument, nonce)
	example.Fatal(err)
	example.Print("cancelled_orders", res.CancelledOrders)
}
