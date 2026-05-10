// Lists every untriggered trigger order on the configured subaccount.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	orders, err := c.GetTriggerOrders(ctx)
	example.Fatal(err)
	example.Print("trigger orders", len(orders))
	for i, o := range orders {
		if i >= 3 {
			break
		}
		example.Print("order", o.OrderID)
		example.Print("  instrument", o.InstrumentName)
		example.Print("  trigger price", o.TriggerPrice.String())
	}
}
