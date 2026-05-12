// Marks specific notifications as seen or hidden over WebSocket.
// Required env: DERIVE_NOTIFICATION_ID.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/amiwrpremium/go-derive/pkg/ws"
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
	raw := os.Getenv("DERIVE_NOTIFICATION_ID")
	if raw == "" {
		log.Fatal("DERIVE_NOTIFICATION_ID required (comma-separated ids)")
	}
	ids := []int64{}
	for _, tok := range strings.Split(raw, ",") {
		id, err := strconv.ParseInt(strings.TrimSpace(tok), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}
	res, err := c.UpdateNotifications(ctx, types.UpdateNotificationsInput{
		NotificationIDs: ids,
		Status:          "seen",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "updated:", res.UpdatedCount)
}
