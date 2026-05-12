// Replaces (cancel + place) one outstanding order in a single round trip.
//
// Replace is the maker-friendly way to re-price without racing the engine.
// Set DERIVE_ORDER_ID to the order you want to cancel; the replacement
// uses the configured BaseAsset / Instrument.
//
// Double-gate the live submission behind DERIVE_RUN_LIVE_ORDERS=1.
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
	id := os.Getenv("DERIVE_ORDER_ID")
	if id == "" {
		log.Fatal("DERIVE_ORDER_ID required")
	}
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually submit")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.Replace(ctx, types.ReplaceOrderInput{
		PlaceOrderInput: types.PlaceOrderInput{
			InstrumentName: instrument,
			Asset:          types.Address(baseAsset),
			Direction:      enums.DirectionBuy,
			OrderType:      enums.OrderTypeLimit,
			TimeInForce:    enums.TimeInForceGTC,
			Amount:         types.MustDecimal("0.01"),
			LimitPrice:     types.MustDecimal("50000"),
			MaxFee:         types.MustDecimal("10"),
		},
		OrderIDToCancel: id,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "cancelled_order:", res.CancelledOrder.OrderID)
	if res.Order != nil {
		fmt.Printf("%-30s %v\n", "new_order:", res.Order.OrderID)
	}
	if res.CreateOrderError != nil {
		fmt.Printf("%-30s %v\n", "create_order_error code:", res.CreateOrderError.Code)
		fmt.Printf("%-30s %v\n", "create_order_error message:", res.CreateOrderError.Message)
	}
	fmt.Printf("%-30s %v\n", "trade count:", len(res.Trades))
}
