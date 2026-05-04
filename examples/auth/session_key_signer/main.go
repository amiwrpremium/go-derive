// Builds an auth.SessionKeySigner — owner address differs from the session key address.
package main

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	key := os.Getenv("DERIVE_SESSION_KEY")
	owner := os.Getenv("DERIVE_OWNER")
	if key == "" || owner == "" {
		log.Fatal("DERIVE_SESSION_KEY and DERIVE_OWNER required")
	}
	s, err := auth.NewSessionKeySigner(key, common.HexToAddress(owner))
	example.Fatal(err)

	example.Print("address (session key)", s.Address())
	example.Print("owner (smart account)", s.Owner())
}
