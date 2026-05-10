// Lists every root wallet (top-of-tree referrer) the engine tracks.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetTreeRoots(ctx)
	example.Fatal(err)
	example.Print("roots (raw bytes)", len(res.Roots))
}
