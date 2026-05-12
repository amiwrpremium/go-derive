// Signs an EIP-712 ActionData and prints the resulting signature.
//
// Note: production code rarely calls SignAction directly — the SDK
// signs order/quote payloads internally inside the client. This
// demo reaches into internal/netconf for the EIP-712 domain so the
// signature can be reproduced standalone.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/internal/netconf"
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

	subStr := os.Getenv("DERIVE_SUBACCOUNT")
	if subStr == "" {
		log.Fatal("DERIVE_SUBACCOUNT required")
	}
	subaccount, err := strconv.ParseInt(subStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_SUBACCOUNT=%q: %v", subStr, err)
	}

	cfg := netconf.Testnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		cfg = netconf.Mainnet()
	}
	domain := cfg.EIP712Domain()

	action := auth.ActionData{
		SubaccountID: subaccount,
		Nonce:        auth.NewNonceGen().Next(),
		Expiry:       1_700_000_000,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sig, err := s.SignAction(ctx, domain, action)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "signature:", sig.Hex())
}
