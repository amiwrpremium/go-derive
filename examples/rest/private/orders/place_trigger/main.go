// Places one stop-loss trigger sell 5% below mark. The order is saved
// server-side in `untriggered` state and only enters the book when the
// mark price crosses the trigger.
//
// Requires: DERIVE_BASE_ASSET, DERIVE_RUN_LIVE_ORDERS=1.
package main

import (
	"log"
	"os"

	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually place an order")
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	tk, err := c.GetTicker(ctx, example.Instrument())
	example.Fatal(err)
	mark := tk.MarkPrice.Inner()
	trigger := mark.Mul(decimal.RequireFromString("0.95"))
	triggerPrice, _ := types.NewDecimal(trigger.String())
	// The resulting market-order limit price is also below mark — the
	// engine still requires a limit_price even for stops; it acts as
	// the slippage cap on the triggered order.
	limit := mark.Mul(decimal.RequireFromString("0.94"))
	limitPrice, _ := types.NewDecimal(limit.String())

	o, err := c.PlaceTriggerOrder(ctx, types.TriggerOrderInput{
		PlaceOrderInput: types.PlaceOrderInput{
			InstrumentName: example.Instrument(),
			Asset:          types.Address(example.BaseAsset()),
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
	example.Fatal(err)
	example.Print("placed", o.OrderID)
	example.Print("status", o.OrderStatus)
	example.Print("trigger_price", o.TriggerPrice.String())
}
