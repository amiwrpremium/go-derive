// Package enums declares the named-string enums used across the SDK.
//
// This file holds [TriggerType] and [TriggerPriceType], the
// trigger-flavour and price-source enums submitted to
// `private/trigger_order`.
package enums

// TriggerType is the flavour of a trigger order — stop-loss versus
// take-profit. Direction has to be set consistently (e.g. a stop-loss
// on a long position is a sell with the trigger below mark).
type TriggerType string

const (
	// TriggerTypeStopLoss fires when the watched price crosses the
	// trigger in the loss-inducing direction.
	TriggerTypeStopLoss TriggerType = "stoploss"
	// TriggerTypeTakeProfit fires when the watched price crosses the
	// trigger in the profit-taking direction.
	TriggerTypeTakeProfit TriggerType = "takeprofit"
)

// Valid reports whether the receiver is one of the defined trigger
// flavours.
func (t TriggerType) Valid() bool {
	switch t {
	case TriggerTypeStopLoss, TriggerTypeTakeProfit:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (t TriggerType) Validate() error {
	if t.Valid() {
		return nil
	}
	return invalid("TriggerType", string(t))
}

// TriggerPriceType is which price the trigger watches — the mark
// price or the index price.
type TriggerPriceType string

const (
	// TriggerPriceTypeMark watches the mark price.
	TriggerPriceTypeMark TriggerPriceType = "mark"
	// TriggerPriceTypeIndex watches the index price. Per the docs the
	// matching engine does not yet support combining "mark"-priced
	// trigger orders with this price source.
	TriggerPriceTypeIndex TriggerPriceType = "index"
)

// Valid reports whether the receiver is one of the defined price
// sources.
func (t TriggerPriceType) Valid() bool {
	switch t {
	case TriggerPriceTypeMark, TriggerPriceTypeIndex:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (t TriggerPriceType) Validate() error {
	if t.Valid() {
		return nil
	}
	return invalid("TriggerPriceType", string(t))
}
