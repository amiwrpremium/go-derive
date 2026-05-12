// Lists per-subaccount portfolio snapshots for every subaccount the
// configured signer's wallet owns. Each snapshot carries the full
// margin breakdown plus collateral / position / open-order arrays.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/ethereum/go-ethereum/common"
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

	portfolios, err := c.GetAllPortfolios(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "portfolio count:", len(portfolios))
	for i, p := range portfolios {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "subaccount:", p.SubaccountID)
		fmt.Printf("%-30s %v\n", "  margin_type:", p.MarginType)
		fmt.Printf("%-30s %v\n", "  subaccount_value:", p.SubaccountValue.String())
		fmt.Printf("%-30s %v\n", "  positions:", len(p.Positions))
		fmt.Printf("%-30s %v\n", "  open_orders:", len(p.OpenOrders))
	}
}
