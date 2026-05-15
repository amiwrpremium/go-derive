// Simulates a margin calculation against the configured subaccount
// and reports pre/post initial- and maintenance-margin values.
// No simulated trade — calling without params reports the current
// subaccount's margin snapshot.
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
	"github.com/amiwrpremium/go-derive/pkg/rest"
)

func main() {
	subStr := os.Getenv("DERIVE_SUBACCOUNT")
	if subStr == "" {
		log.Fatal("DERIVE_SUBACCOUNT required")
	}
	subaccount, err := strconv.ParseInt(subStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_SUBACCOUNT=%q: %v", subStr, err)
	}
	key := os.Getenv("DERIVE_SESSION_KEY")
	if key == "" {
		log.Fatal("DERIVE_SESSION_KEY required")
	}
	var signer auth.Signer
	if owner := os.Getenv("DERIVE_OWNER"); owner != "" {
		signer, err = auth.NewSessionKeySigner(key, common.HexToAddress(owner))
	} else {
		signer, err = auth.NewLocalSigner(key)
	}
	if err != nil {
		log.Fatalf("signer: %v", err)
	}

	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork, rest.WithSigner(signer), rest.WithSubaccount(subaccount))
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.GetMargin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "subaccount_id:", res.SubaccountID)
	fmt.Printf("%-30s %v\n", "pre_initial_margin:", res.PreInitialMargin.String())
	fmt.Printf("%-30s %v\n", "post_initial_margin:", res.PostInitialMargin.String())
	fmt.Printf("%-30s %v\n", "is_valid_trade:", res.IsValidTrade)
}
