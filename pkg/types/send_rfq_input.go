// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the user-facing input DTO for the `SendRFQ` method,
// wrapping `private/send_rfq`.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// SendRFQInput parameterises a request-for-quote submission. Only
// [SendRFQInput.Legs] is required; every other field is optional and
// AND-filters against the matching makers / accepted quotes.
//
// Wire keys mirror docs.derive.xyz/reference/post_private-send-rfq:
// max_total_cost / min_total_cost / partial_fill_step / etc. The SDK
// fills in subaccount id from the configured client.
type SendRFQInput struct {
	// Legs is the per-instrument shape of the requested trade.
	// Required and non-empty.
	Legs []RFQLeg
	// Counterparties, when non-empty, restricts the RFQ to this
	// whitelist of maker wallet addresses. Empty means the RFQ
	// is broadcast to every eligible maker.
	Counterparties []Address
	// PreferredDirection is an optional hint to makers about the
	// direction the taker is more likely to fill.
	PreferredDirection enums.Direction
	// ReducingDirection is an optional risk-reducing marker — set
	// to indicate the RFQ would reduce the taker's existing
	// position. Some maker programs price these tighter.
	ReducingDirection enums.Direction
	// Label is an optional user-defined tag attached to the RFQ.
	Label string
	// MaxTotalCost caps the taker-side total cost the engine will
	// accept on an executed quote. Zero defers to the engine's
	// default (no upper cap).
	MaxTotalCost Decimal
	// MinTotalCost floors the taker-side total cost the engine
	// will accept on an executed quote. Zero defers to the
	// engine's default (no lower floor).
	MinTotalCost Decimal
	// PartialFillStep is the minimum increment for partial fills.
	// Zero defers to the engine's default (no partial fills).
	PartialFillStep Decimal
	// Client is an optional caller-defined tag echoed back on the
	// RFQ record.
	Client string
	// ReferralCode is an optional referral code applied to fills.
	ReferralCode string
	// ExtraFee is an optional caller-paid tip on top of the
	// standard fee schedule. Denominated in quote currency.
	ExtraFee Decimal
}
