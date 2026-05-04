// Fetches Derive's server time over REST.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	t, err := c.GetTime(ctx)
	example.Fatal(err)
	example.Print("server time (ms)", t)
}
