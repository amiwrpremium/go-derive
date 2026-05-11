// Paginates the configured subaccount's quotes — own quotes either
// active or in any historical state.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	quotes, page, err := c.GetQuotes(ctx, map[string]any{"page_size": 10})
	example.Fatal(err)
	example.Print("count", len(quotes))
	example.Print("total pages", page.NumPages)
}
