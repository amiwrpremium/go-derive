// Builds the facade with private credentials and uses both REST and WS.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/derive"
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

	network := derive.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		network = derive.WithMainnet()
	}
	c, err := derive.NewClient(network, derive.WithSigner(s), derive.WithSubaccount(subaccount))
	if err != nil {
		log.Fatalf("derive.NewClient: %v", err)
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := c.WS.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	if err := c.WS.Login(ctx); err != nil {
		log.Fatal(err)
	}

	sa, err := c.REST.GetSubaccount(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "subaccount:", sa.SubaccountID)
	fmt.Printf("%-30s %v\n", "ws connected:", c.WS.IsConnected())
}
