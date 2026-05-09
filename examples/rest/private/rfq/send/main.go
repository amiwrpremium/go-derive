// Sends a one-leg RFQ.
package main

import "github.com/amiwrpremium/go-derive"

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	rfq, err := c.SendRFQ(ctx, []derive.RFQLeg{
		{InstrumentName: example.Instrument(), Direction: derive.DirectionBuy, Amount: derive.MustDecimal("0.1")},
	}, derive.MustDecimal("10"))
	example.Fatal(err)
	example.Print("rfq id", rfq.RFQID)
	example.Print("status", rfq.Status)
}
