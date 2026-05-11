// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the user-facing input DTO for
// `private/rfq_get_best_quote`, the taker-side pre-flight that
// returns the best executable price for an RFQ shape without
// actually submitting one.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// BestQuoteInput is the typed body for `private/rfq_get_best_quote`.
// It returns the best executable price the engine can match
// against the supplied RFQ shape, plus margin-impact estimates,
// without entering the actual RFQ flow.
//
// No signing is required — the call is a pure lookup.
type BestQuoteInput struct {
	// SubaccountID is the taker subaccount. Zero defaults to the
	// client-configured subaccount.
	SubaccountID int64
	// Legs describe the unpriced RFQ shape.
	Legs []RFQLeg
	// Direction is the side the taker would execute on. Defaults to
	// buy when unset.
	Direction enums.Direction
	// PreferredDirection is the maker side the engine should prefer
	// when both sides are available.
	PreferredDirection enums.Direction
	// Counterparties optionally restricts the lookup to this list of
	// maker account addresses.
	Counterparties []string
	// Label is an optional user-defined tag.
	Label string
	// Client is an optional client identifier.
	Client string
	// ExtraFee is an optional USDC tip (0.000001-1000) the taker
	// adds on top of the engine fee.
	ExtraFee Decimal
	// MaxTotalCost is the maximum total cost the taker will accept;
	// the engine treats it as confidential.
	MaxTotalCost Decimal
	// MinTotalCost is the minimum total cost the taker will accept;
	// confidential to the engine.
	MinTotalCost Decimal
	// PartialFillStep is the step size in base units the engine may
	// use for partial fills.
	PartialFillStep Decimal
	// ReferralCode is an optional referral identifier.
	ReferralCode string
	// RFQID is the optional RFQ id to retrieve a best quote against
	// an already-broadcast RFQ.
	RFQID string
}

// Validate performs schema-level checks on the receiver. Returns
// nil on success or a wrapped [ErrInvalidParams].
func (in BestQuoteInput) Validate() error {
	if len(in.Legs) == 0 {
		return invalidParam("legs", "must have at least one leg")
	}
	for i := range in.Legs {
		if err := in.Legs[i].Validate(); err != nil {
			return err
		}
	}
	if in.Direction != "" {
		if err := in.Direction.Validate(); err != nil {
			return invalidParam("direction", err.Error())
		}
	}
	if in.PreferredDirection != "" {
		if err := in.PreferredDirection.Validate(); err != nil {
			return invalidParam("preferred_direction", err.Error())
		}
	}
	return nil
}
