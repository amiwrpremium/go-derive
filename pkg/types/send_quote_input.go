// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the user-facing input DTOs for the maker
// quote-send and quote-replace flows.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// SendQuoteInput is the typed body for `private/send_quote`. A maker
// uses this to respond to an open RFQ with a multi-leg priced
// quote.
//
// Quote signing is not yet handled by the SDK: callers must
// pre-compute [Signature], [Signer], [SignatureExpirySec] and
// [Nonce] using their own signing flow before calling
// [methods.API.SendQuote]. The remaining fields are forwarded to
// the engine verbatim.
type SendQuoteInput struct {
	// RFQID is the open RFQ this quote responds to.
	RFQID string
	// SubaccountID is the maker subaccount issuing the quote. When
	// zero the SDK fills it in from the client configuration.
	SubaccountID int64
	// Direction is the quote's overall direction. `buy` means each
	// leg trades in its own direction; `sell` means each leg trades
	// in the opposite direction.
	Direction enums.Direction
	// Legs are the priced quote legs. At least one leg is required.
	Legs []QuoteLeg
	// MaxFee is the maximum dollar fee for the full trade. The
	// engine rejects the quote if the estimated fee exceeds it.
	MaxFee Decimal
	// Nonce is the unique per-quote nonce; the engine deduplicates
	// against this value.
	Nonce uint64
	// Signature is the caller-computed Ethereum signature over the
	// quote payload.
	Signature string
	// Signer is the wallet or session-key address that produced
	// [Signature].
	Signer string
	// SignatureExpirySec is the Unix timestamp (seconds) after which
	// the signature is no longer valid. Must be at least 310 seconds
	// in the future per the engine; the quote expires at
	// SignatureExpirySec - 300.
	SignatureExpirySec int64
	// Label is an optional user-defined tag attached to the quote.
	Label string
	// MMP marks the quote as eligible for market-maker protection
	// cancellations.
	MMP bool
	// Client is the optional client identifier echoed back in
	// quote-status notifications.
	Client string
}

// Validate performs schema-level checks on the receiver. Returns
// nil on success or a wrapped [ErrInvalidParams].
func (in SendQuoteInput) Validate() error {
	if in.RFQID == "" {
		return invalidParam("rfq_id", "required")
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
	}
	if in.MaxFee.Sign() < 0 {
		return invalidParam("max_fee", "must be non-negative")
	}
	if in.Nonce == 0 {
		return invalidParam("nonce", "required")
	}
	if in.Signature == "" {
		return invalidParam("signature", "required")
	}
	if in.Signer == "" {
		return invalidParam("signer", "required")
	}
	if in.SignatureExpirySec == 0 {
		return invalidParam("signature_expiry_sec", "required")
	}
	return nil
}

// ReplaceQuoteInput is the typed body for `private/replace_quote`.
// It cancels one outstanding maker quote and submits a replacement
// in a single round trip.
//
// Exactly one of [QuoteIDToCancel] / [NonceToCancel] must be set.
// All other fields are inherited from the embedded
// [SendQuoteInput] and signed exactly like a fresh send-quote.
type ReplaceQuoteInput struct {
	SendQuoteInput
	// QuoteIDToCancel identifies the quote to cancel by its
	// engine-assigned id. Set this when you have the id from a prior
	// send_quote response.
	QuoteIDToCancel string
	// NonceToCancel identifies the quote to cancel by its signed
	// nonce. Set this when cancelling a quote whose id has not yet
	// been received.
	NonceToCancel uint64
}

// Validate performs schema-level checks on the receiver. Wraps the
// embedded [SendQuoteInput.Validate] with the cancel-target
// requirement.
func (in ReplaceQuoteInput) Validate() error {
	if err := in.SendQuoteInput.Validate(); err != nil {
		return err
	}
	hasID := in.QuoteIDToCancel != ""
	hasNonce := in.NonceToCancel != 0
	if !hasID && !hasNonce {
		return invalidParam("quote_id_to_cancel|nonce_to_cancel", "one of quote_id_to_cancel or nonce_to_cancel is required")
	}
	if hasID && hasNonce {
		return invalidParam("quote_id_to_cancel|nonce_to_cancel", "must not set both quote_id_to_cancel and nonce_to_cancel")
	}
	return nil
}
