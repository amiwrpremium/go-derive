// Lists active maker incentive programs. Each program has its own
// epoch, the asset types and currencies it covers, and the rewards
// paid out.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	programs, err := c.GetMakerPrograms(ctx)
	example.Fatal(err)
	example.Print("program count", len(programs))
	for i, p := range programs {
		if i >= 5 {
			break
		}
		example.Print("program", p.Name)
		example.Print("  asset_types", p.AssetTypes)
		example.Print("  currencies", p.Currencies)
		example.Print("  min_notional", p.MinNotional.String())
	}
}
