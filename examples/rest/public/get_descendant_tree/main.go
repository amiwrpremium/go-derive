// Returns the referee tree rooted at one wallet (or invite code).
// Required env: DERIVE_WALLET_OR_INVITE_CODE.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork)
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	root := os.Getenv("DERIVE_WALLET_OR_INVITE_CODE")
	if root == "" {
		log.Fatal("DERIVE_WALLET_OR_INVITE_CODE required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tree, err := c.GetDescendantTree(ctx, types.DescendantTreeQuery{WalletOrInviteCode: root})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "parent:", tree.Parent)
	fmt.Printf("%-30s %v\n", "descendants (raw bytes):", len(tree.Descendants))
}
