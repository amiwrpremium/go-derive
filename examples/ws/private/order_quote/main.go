// Pre-flights a hypothetical order for the configured subaccount
// over WebSocket. Private counterpart to public/order_quote.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.OrderQuote(ctx, types.PlaceOrderInput{
		InstrumentName: example.Instrument(),
		Asset:          types.Address(example.BaseAsset()),
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		TimeInForce:    enums.TimeInForceGTC,
		Amount:         types.MustDecimal("0.01"),
		LimitPrice:     types.MustDecimal("50000"),
		MaxFee:         types.MustDecimal("10"),
	})
	example.Fatal(err)
	example.Print("is_valid", res.IsValid)
	example.Print("estimated_fee", res.EstimatedFee.String())
}
