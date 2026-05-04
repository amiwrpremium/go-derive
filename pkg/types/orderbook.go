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
