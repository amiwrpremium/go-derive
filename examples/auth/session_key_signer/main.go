// Builds an derive.SessionKeySigner — owner address differs from the session key address.
package main

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	key := os.Getenv("DERIVE_SESSION_KEY")
	owner := os.Getenv("DERIVE_OWNER")
	if key == "" || owner == "" {
		log.Fatal("DERIVE_SESSION_KEY and DERIVE_OWNER required")
	}
	s, err := derive.NewSessionKeySigner(key, common.HexToAddress(owner))
	example.Fatal(err)

	example.Print("address (session key)", s.Address())
	example.Print("owner (smart account)", s.Owner())
}
