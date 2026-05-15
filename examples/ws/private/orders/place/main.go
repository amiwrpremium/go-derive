// Places one limit order over WebSocket. Requires DERIVE_BASE_ASSET and
// DERIVE_RUN_LIVE_ORDERS=1.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	instrument := os.Getenv("DERIVE_INSTRUMENT")
	if instrument == "" {
		instrument = "BTC-PERP"
	}
	baseAsset := common.HexToAddress(os.Getenv("DERIVE_BASE_ASSET"))
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
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1")
	}
	tk, err := c.GetTicker(ctx, types.TickerQuery{Name: instrument})
	if err != nil {
		log.Fatal(err)
	}
	limit := tk.MarkPrice.Inner().Mul(decimal.RequireFromString("0.95"))
	price, _ := types.NewDecimal(limit.String())

	o, trades, err := c.PlaceOrder(ctx, types.PlaceOrderInput{
		InstrumentName: instrument,
		Asset:          types.Address(baseAsset),
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		TimeInForce:    enums.TimeInForceGTC,
		Amount:         types.MustDecimal("0.001"),
		LimitPrice:     price,
		MaxFee:         types.MustDecimal("10"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "placed:", o.OrderID)
	fmt.Printf("%-30s %v\n", "trades on submission:", len(trades))
}
