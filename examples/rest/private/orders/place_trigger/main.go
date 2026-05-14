// Places one stop-loss trigger sell 5% below mark. The order is saved
// server-side in `untriggered` state and only enters the book when the
// mark price crosses the trigger.
//
// Requires: DERIVE_BASE_ASSET, DERIVE_RUN_LIVE_ORDERS=1.
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
	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/types"
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

	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork, rest.WithSigner(signer), rest.WithSubaccount(subaccount))
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually place an order")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tk, err := c.GetTicker(ctx, instrument)
	if err != nil {
		log.Fatal(err)
	}
	mark := tk.MarkPrice.Inner()
	trigger := mark.Mul(decimal.RequireFromString("0.95"))
	triggerPrice, _ := types.NewDecimal(trigger.String())
	// The resulting market-order limit price is also below mark — the
	// engine still requires a limit_price even for stops; it acts as
	// the slippage cap on the triggered order.
	limit := mark.Mul(decimal.RequireFromString("0.94"))
	limitPrice, _ := types.NewDecimal(limit.String())

	o, _, err := c.PlaceTriggerOrder(ctx, types.TriggerOrderInput{
		PlaceOrderInput: types.PlaceOrderInput{
			InstrumentName: instrument,
			Asset:          types.Address(baseAsset),
			Direction:      enums.DirectionSell,
			OrderType:      enums.OrderTypeLimit,
			TimeInForce:    enums.TimeInForceGTC,
			Amount:         types.MustDecimal("0.001"),
			LimitPrice:     limitPrice,
			MaxFee:         types.MustDecimal("10"),
			Label:          "go-derive-trigger",
		},
		TriggerType:      enums.TriggerTypeStopLoss,
		TriggerPriceType: enums.TriggerPriceTypeMark,
		TriggerPrice:     triggerPrice,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "placed:", o.OrderID)
	fmt.Printf("%-30s %v\n", "status:", o.OrderStatus)
	fmt.Printf("%-30s %v\n", "trigger_price:", o.TriggerPrice.String())
}
