// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// All numeric fields use [Decimal], a thin wrapper around shopspring/decimal,
// so price/size/fee values never lose precision through float64 round-trips.
// On the wire, [Decimal] reads and writes JSON strings (Derive's preferred
// representation); a fallback path also accepts JSON numbers for resilience.
//
// Identifier types ([Address], [TxHash], [MillisTime]) carry the same
// round-trip guarantees: each one preserves the canonical wire format
// regardless of how Go marshals the surrounding struct.
//
// # Why named types
//
// Plain string and int64 fields would parse just fine, but named types let
// the SDK enforce invariants at construction time (NewAddress checksum
// check, NewDecimal precision check) and let callers tell at a glance which
// values are amounts vs prices vs subaccount ids.
package types

import (
	"encoding/json"
	"fmt"
)

// OrderBookLevel is one [price, amount] pair on either side of an order book.
//
// Derive serializes order-book levels as two-element JSON arrays rather than
// objects, e.g. ["65000", "1.5"]. A custom Marshal/Unmarshal preserves that
// wire format while exposing readable Price / Amount fields at the call site.
type OrderBookLevel struct {
	// Price is the resting limit price.
	Price Decimal
	// Amount is the resting size at that price.
	Amount Decimal
}

// MarshalJSON encodes the level as a [price, amount] JSON array.
func (l OrderBookLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]Decimal{l.Price, l.Amount})
}

// UnmarshalJSON decodes a [price, amount] JSON array into the receiver.
// Object-shaped input returns an error.
func (l *OrderBookLevel) UnmarshalJSON(b []byte) error {
	var arr [2]Decimal
	if err := json.Unmarshal(b, &arr); err != nil {
		return fmt.Errorf("types: orderbook level: %w", err)
	}
	l.Price = arr[0]
	l.Amount = arr[1]
	return nil
}

// OrderBook is a snapshot of the order book at a point in time.
//
// Both sides are sorted: [Bids] is descending by price, [Asks] is ascending.
type OrderBook struct {
	// InstrumentName identifies the market this snapshot belongs to.
	InstrumentName string `json:"instrument_name"`
	// PublishID is the monotonically incrementing publish counter. Use it
	// to detect dropped or out-of-order messages on the orderbook channel.
	PublishID int64 `json:"publish_id,omitempty"`
	// Bids is the buy side of the book.
	Bids []OrderBookLevel `json:"bids"`
	// Asks is the sell side of the book.
	Asks []OrderBookLevel `json:"asks"`
	// Timestamp is the engine-side capture time.
	Timestamp MillisTime `json:"timestamp"`
	// PublishTime is the time the snapshot was published over the wire.
	// Compare against Timestamp to gauge engine-to-client latency.
	PublishTime MillisTime `json:"publish_time,omitempty"`
}
