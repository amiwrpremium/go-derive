// Sends a one-leg RFQ over WebSocket.
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

	rfq, err := c.SendRFQ(ctx, []types.RFQLeg{
		{InstrumentName: example.Instrument(), Direction: enums.DirectionBuy, Amount: types.MustDecimal("0.1")},
	}, types.MustDecimal("10"))
	example.Fatal(err)
	example.Print("rfq id", rfq.RFQID)
}
