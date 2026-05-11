// Pre-flights a hypothetical order through the matching engine for
// the configured subaccount. Private counterpart to
// public/order_quote — same response shape, but the engine accounts
// for the subaccount's actual margin / collateral.
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
	example.Print("estimated_fill_price", res.EstimatedFillPrice.String())
	example.Print("max_amount", res.MaxAmount.String())
}
