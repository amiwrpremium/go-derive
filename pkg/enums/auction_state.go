// Package enums declares the named-string enums used across the SDK.
//
// This file holds [AuctionState], the lifecycle state reported on the
// `auctions.watch` WebSocket channel.
package enums

// AuctionState is the lifecycle state reported on the `auctions.watch`
// WebSocket channel.
type AuctionState string

const (
	// AuctionStateOngoing — the auction is in progress and accepting bids.
	AuctionStateOngoing AuctionState = "ongoing"
	// AuctionStateEnded — the auction has concluded.
	AuctionStateEnded AuctionState = "ended"
)

// Valid reports whether the receiver is one of the defined auction states.
func (s AuctionState) Valid() bool {
	switch s {
	case AuctionStateOngoing, AuctionStateEnded:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (s AuctionState) Validate() error {
	if s.Valid() {
		return nil
	}
	return invalid("AuctionState", string(s))
}
