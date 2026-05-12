// Exports the configured subaccount's expired and cancelled orders.
// Returns presigned S3 URLs you can fetch the CSV data from.
//
// Required env: DERIVE_WALLET (the wallet to scope the export to).
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
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.ExpiredAndCancelledHistory(ctx, types.ExpiredAndCancelledHistoryInput{
		Wallet:    wallet,
		ExpirySec: 3600,
	})
	example.Fatal(err)
	example.Print("urls", len(res.PresignedURLs))
	for i, u := range res.PresignedURLs {
		if i >= 3 {
			break
		}
		example.Print("  url", u)
	}
}
