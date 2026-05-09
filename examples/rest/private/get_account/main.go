// Fetches wallet-level account information for the configured signer.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	acc, err := c.GetAccount(ctx)
	example.Fatal(err)
	example.Print("wallet", acc.Wallet)
	example.Print("subaccount count", len(acc.SubaccountIDs))
	example.Print("perp taker fee", acc.FeeInfo.PerpTakerFee.String())
}
