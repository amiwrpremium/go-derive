// Pre-flights a hypothetical order through the matching engine
// via public/order_quote. The endpoint requires a signed body even
// though it is not subscription-gated, so this uses the private
// client to populate the signing fields.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.OrderQuotePublic(ctx, types.PlaceOrderInput{
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
	example.Print("estimated_fill_price", res.EstimatedFillPrice.String())
	example.Print("estimated_order_status", string(res.EstimatedOrderStatus))
	example.Print("max_amount", res.MaxAmount.String())
}
