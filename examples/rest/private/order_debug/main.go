// Previews an order without submitting it. Returns the engine's view of
// fees and margin impact, useful for sanity-checking signed payloads in
// CI or pre-flight.
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

	dbg, err := c.OrderDebug(ctx, types.PlaceOrderInput{
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
	example.Print("typed_data_hash", dbg.TypedDataHash)
	example.Print("action_hash", dbg.ActionHash)
	example.Print("module", dbg.RawData.Module)
}
