// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the input DTO for the two batch-cancel methods
// `private/cancel_batch_quotes` and `private/cancel_batch_rfqs`.
package types

// CancelBatchInput filters which RFQs or quotes to cancel. All
// filters AND together; if every field is its zero value the call
// is a no-op (cancels nothing).
//
// Used by both [methods.API.CancelBatchQuotes] and
// [methods.API.CancelBatchRFQs]. The QuoteID filter is only
// honoured by CancelBatchQuotes.
type CancelBatchInput struct {
	// SubaccountID restricts the cancel to one subaccount. Zero
	// defaults to the client-configured subaccount.
	SubaccountID int64
	// RFQID restricts the cancel to quotes/RFQs for this RFQ.
	RFQID string
	// QuoteID restricts the cancel to one specific quote
	// (CancelBatchQuotes only).
	QuoteID string
	// Label restricts the cancel to RFQs/quotes that were tagged
	// with this label at submission time.
	Label string
	// Nonce restricts the cancel to one specific signed nonce.
	Nonce uint64
}

// Validate enforces that at least one filter beyond SubaccountID is
// populated. An all-zero input would no-op server-side, almost always
// meaning the caller forgot to set a filter; this catches that case
// before a wasted round-trip.
func (in CancelBatchInput) Validate() error {
	if in.RFQID == "" && in.QuoteID == "" && in.Label == "" && in.Nonce == 0 {
		return invalidParam("filter",
			"at least one of rfq_id, quote_id, label, nonce is required")
	}
	return nil
}
