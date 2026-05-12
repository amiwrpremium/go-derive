// Builds an auth.LocalSigner from the configured private key.
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
	if key == "" {
		log.Fatal("DERIVE_SESSION_KEY required")
	}
	var s auth.Signer
	var err error
	if owner := os.Getenv("DERIVE_OWNER"); owner != "" {
		s, err = auth.NewSessionKeySigner(key, common.HexToAddress(owner))
	} else {
		s, err = auth.NewLocalSigner(key)
	}
	if err != nil {
		log.Fatalf("signer: %v", err)
	}
	fmt.Printf("%-30s %v\n", "address:", s.Address())
	fmt.Printf("%-30s %v\n", "owner:", s.Owner())
}
