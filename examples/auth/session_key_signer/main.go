// Builds an auth.SessionKeySigner — owner address differs from the session key address.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	key := os.Getenv("DERIVE_SESSION_KEY")
	owner := os.Getenv("DERIVE_OWNER")
	if key == "" || owner == "" {
		log.Fatal("DERIVE_SESSION_KEY and DERIVE_OWNER required")
	}
	s, err := auth.NewSessionKeySigner(key, common.HexToAddress(owner))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-30s %v\n", "address (session key):", s.SessionAddress())
	fmt.Printf("%-30s %v\n", "owner (smart account):", s.OwnerAddress())
}
