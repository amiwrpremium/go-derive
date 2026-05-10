package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// GetOrdersFilter narrows a paginated `private/get_orders` request.
// Each field is optional; the zero value asks the engine for
// unfiltered results.
type GetOrdersFilter struct {
	// InstrumentName filters to one instrument.
	InstrumentName string
	// Label filters to orders carrying the user-defined label.
	Label string
	// Status filters by order status. The wire enum is
	// open / filled / cancelled / expired / untriggered / algo_active
	// (see [enums.OrderStatus]).
	Status enums.OrderStatus
}
