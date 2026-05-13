// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the user-facing input DTO for
// `private/execute_quote`, the taker side of the RFQ workflow.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// ExecuteQuoteInput is the typed body for `private/execute_quote`.
// A taker uses this to settle against one specific maker quote
// that came back on an open RFQ.
//
// The SDK signs the per-execute EIP-712 payload internally. Each
// leg must carry both the engine-facing fields and the on-chain
// identifiers (`Asset`, `SubID`) used by the RFQ module's hash;
// obtain them via `public/get_instrument`. The SDK internally
// inverts the global direction when computing the per-leg signed
// amount, since the taker takes the opposite side of the maker
// quote.
type ExecuteQuoteInput struct {
	// RFQID is the RFQ the taker initially sent.
	RFQID string
	// QuoteID is the maker quote being executed against.
	QuoteID string
	// SubaccountID is the taker subaccount. Zero defaults to the
	// client-configured subaccount.
	SubaccountID int64
	// Direction is the taker's intended trade direction; must match
	// the side they want to take on the legs.
	Direction enums.Direction
	// Legs mirror the maker quote's legs and prices; the engine
	// uses them to confirm the taker is settling the quote they
	// expect.
	Legs []QuoteLeg
	// MaxFee is the maximum dollar fee for the full trade.
	MaxFee Decimal
	// Label is an optional user-defined tag.
	Label string
	// EnableTakerProtection turns on Derive's taker-side protection
	// (rejects the execute if the engine sees a worse quote at fill
	// time).
	EnableTakerProtection bool
	// Client is the optional client identifier echoed back in
	// trade notifications.
	Client string
}

// Validate performs schema-level checks on the receiver. Returns
// nil on success or a wrapped [ErrInvalidParams].
func (in ExecuteQuoteInput) Validate() error {
	if in.RFQID == "" {
		return invalidParam("rfq_id", "required")
	}
	if in.QuoteID == "" {
		return invalidParam("quote_id", "required")
	}
	if err := in.Direction.Validate(); err != nil {
		return invalidParam("direction", err.Error())
	}
	if len(in.Legs) == 0 {
		return invalidParam("legs", "must have at least one leg")
	}
	for i := range in.Legs {
		if in.Legs[i].InstrumentName == "" {
			return invalidParam("legs.instrument_name", "required")
		}
		if err := in.Legs[i].Direction.Validate(); err != nil {
			return invalidParam("legs.direction", err.Error())
		}
		if in.Legs[i].Amount.Sign() <= 0 {
			return invalidParam("legs.amount", "must be positive")
		}
		if (in.Legs[i].Asset == Address{}) {
			return invalidParam("legs.asset", "required for SDK-side signing")
		}
	}
	if in.MaxFee.Sign() < 0 {
		return invalidParam("max_fee", "must be non-negative")
	}
	return nil
}
