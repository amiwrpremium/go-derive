// Cancels every order on the configured subaccount that carries the
// given nonce. Required env: DERIVE_INSTRUMENT (e.g. BTC-PERP) and
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

	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.CancelByNonce(ctx, instrument, nonce)
	example.Fatal(err)
	example.Print("cancelled_orders", res.CancelledOrders)
}
