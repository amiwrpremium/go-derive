// Returns the referee tree rooted at one wallet (or invite code).
// Required env: DERIVE_WALLET_OR_INVITE_CODE.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	root := os.Getenv("DERIVE_WALLET_OR_INVITE_CODE")
	if root == "" {
		log.Fatal("DERIVE_WALLET_OR_INVITE_CODE required")
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	tree, err := c.GetDescendantTree(ctx, root)
	example.Fatal(err)
	example.Print("parent", tree.Parent)
	example.Print("descendants (raw bytes)", len(tree.Descendants))
}
