// Signs an EIP-191 timestamp for use as the X-LyraSignature header.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := s.SignAuthHeader(ctx, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "signature:", sig.Hex())
}
