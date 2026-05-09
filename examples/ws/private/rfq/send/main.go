// Sends a one-leg RFQ over WebSocket.
package main

import "github.com/amiwrpremium/go-derive"

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	rfq, err := c.SendRFQ(ctx, []derive.RFQLeg{
		{InstrumentName: example.Instrument(), Direction: derive.DirectionBuy, Amount: derive.MustDecimal("0.1")},
	}, derive.MustDecimal("10"))
	example.Fatal(err)
	example.Print("rfq id", rfq.RFQID)
}
