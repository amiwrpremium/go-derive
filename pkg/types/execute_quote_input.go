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
// Quote signing is not yet handled by the SDK: callers must
// pre-compute [Signature], [Signer], [SignatureExpirySec] and
// [Nonce] from their own signing flow. The remaining fields are
// forwarded to the engine verbatim.
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
	// Nonce is the unique per-execution nonce.
	Nonce uint64
	// Signature is the caller-computed Ethereum signature over the
	// execution payload.
	Signature string
	// Signer is the wallet or session-key address that produced
	// [Signature].
	Signer string
	// SignatureExpirySec is the Unix timestamp (seconds) after which
	// the signature is no longer valid. Must be at least 310 seconds
	// in the future.
	SignatureExpirySec int64
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
