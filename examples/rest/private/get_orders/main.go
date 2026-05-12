// Paginates orders on the configured subaccount, optionally filtered
// by instrument / label / status. Pass nil filter to omit filters
// and page through every order on the subaccount.
//
// Counterpart to GetOrderHistory (time-window-based pagination).
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
	"github.com/amiwrpremium/go-derive/pkg/types"
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

	orders, page, err := c.GetOrders(ctx, types.PageRequest{PageSize: 10}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "orders:", len(orders))
	fmt.Printf("%-30s %v\n", "total pages:", page.NumPages)
	for i, o := range orders {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "order:", o.OrderID)
		fmt.Printf("%-30s %v\n", "  status:", string(o.OrderStatus))
		fmt.Printf("%-30s %v\n", "  instrument:", o.InstrumentName)
	}
}
