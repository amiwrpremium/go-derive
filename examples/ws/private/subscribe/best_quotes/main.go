// Streams the running best-quote state for every open RFQ on the
// configured subaccount. Wire channel: `{subaccount_id}.best.quotes`.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/ws"
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	wsNetwork := ws.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		wsNetwork = ws.WithMainnet()
	}
	c, err := ws.New(wsNetwork, ws.WithSigner(signer), ws.WithSubaccount(subaccount))
	if err != nil {
		log.Fatalf("ws.New: %v", err)
	}
	defer c.Close()
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("ws.Connect: %v", err)
	}
	if err := c.Login(ctx); err != nil {
		log.Fatalf("ws.Login: %v", err)
	}
	sub, err := c.SubscribeBestQuotes(ctx, subaccount)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case batch, ok := <-sub.Updates():
			if !ok {
				return
			}
			fmt.Printf("%-30s %v\n", "events in batch:", len(batch))
			for i, ev := range batch {
				if i >= 3 {
					break
				}
				fmt.Printf("%-30s %v\n", "rfq_id:", ev.RFQID)
				if ev.Error != nil {
					fmt.Printf("%-30s %v\n", "  error:", ev.Error.Message)
				} else if ev.Result != nil {
					fmt.Printf("%-30s %v\n", "  is_valid:", ev.Result.IsValid)
				}
			}
		}
	}
}
