// Exports the configured subaccount's expired and cancelled orders.
// Returns presigned S3 URLs you can fetch the CSV data from.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.ExpiredAndCancelledHistory(ctx, nil)
	example.Fatal(err)
	example.Print("urls", len(res.PresignedURLs))
	for i, u := range res.PresignedURLs {
		if i >= 3 {
			break
		}
		example.Print("  url", u)
	}
}
