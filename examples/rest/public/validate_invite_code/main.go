// Validates one invite code. Required env: DERIVE_INVITE_CODE.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
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
	code := os.Getenv("DERIVE_INVITE_CODE")
	if code == "" {
		log.Fatal("DERIVE_INVITE_CODE required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, err := c.ValidateInviteCode(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "status:", status)
}
