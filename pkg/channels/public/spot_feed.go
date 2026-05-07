// Package public — see orderbook.go for the overview.
package public

import (
	"encoding/json"
	"fmt"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// SpotFeed subscribes to oracle price-feed updates for one currency.
//
// The dotted server-side channel name is:
//
//	spot_feed.{currency}
//
// Currency is the underlying asset symbol (e.g. "BTC", "ETH"). The
// payload includes the current oracle price plus the 24-hour-prior
// reading, so consumers can render a delta without an extra round trip.
//
// Pair this descriptor with T = [types.SpotFeed].
type SpotFeed struct {
	// Currency is the underlying asset symbol (e.g. "BTC").
	Currency string
}

// Name returns the dotted server-side channel string.
func (s SpotFeed) Name() string { return fmt.Sprintf("spot_feed.%s", s.Currency) }

// Decode parses an inbound notification payload into a [types.SpotFeed].
func (SpotFeed) Decode(raw json.RawMessage) (any, error) {
	var sf types.SpotFeed
	if err := json.Unmarshal(raw, &sf); err != nil {
		return nil, err
	}
	return sf, nil
}
