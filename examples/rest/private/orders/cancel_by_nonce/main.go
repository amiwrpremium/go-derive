// Cancels every order on the configured subaccount that carries the
// given nonce. Required env: DERIVE_INSTRUMENT (e.g. BTC-PERP) and
// DERIVE_NONCE.
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
	instrument := os.Getenv("DERIVE_INSTRUMENT")
	nonceStr := os.Getenv("DERIVE_NONCE")
	if instrument == "" || nonceStr == "" {
		log.Fatal("DERIVE_INSTRUMENT and DERIVE_NONCE required")
	}
	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.CancelByNonce(ctx, types.CancelByNonceInput{InstrumentName: instrument, Nonce: nonce})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "cancelled_orders:", res.CancelledOrders)
}
