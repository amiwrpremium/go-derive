// Lists every active perpetual instrument across all currencies,
// paginated. Counterpart to GetInstruments which filters by
// currency.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	insts, page, err := c.GetAllInstruments(ctx, enums.InstrumentTypePerp, false, types.PageRequest{PageSize: 50})
	example.Fatal(err)
	example.Print("instrument count", len(insts))
	example.Print("total pages", page.NumPages)
	for i, in := range insts {
		if i >= 5 {
			break
		}
		example.Print("instrument", in.Name)
	}
}
