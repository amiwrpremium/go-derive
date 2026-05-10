// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the per-event payload of the
// `{subaccount_id}.best.quotes` WebSocket channel — a running stream
// of best-quote state for each open RFQ on the subaccount.
package types

// BestQuoteFeedEvent is one event from the
// `{subaccount_id}.best.quotes` channel. Each event carries the
// engine's current best-quote state for one RFQ on the subaccount,
// either as a `Result` (when the engine could compute it) or an
// `Error` (when the underlying `private/rfq_get_best_quote` RPC
// call failed).
//
// Mirrors the per-event shape per
// docs.derive.xyz/reference/subaccount_id-best-quotes.
type BestQuoteFeedEvent struct {
	// RFQID identifies the RFQ this event is about.
	RFQID string `json:"rfq_id"`
	// Error is set when the engine could not compute a best quote
	// (the underlying `rfq_get_best_quote` RPC failed). Mutually
	// exclusive with Result.
	Error *RPCError `json:"error,omitempty"`
	// Result is the current best-quote state. Same shape as the
	// return of [API.RFQGetBestQuote] (the
	// `private/rfq_get_best_quote` RPC). Nil when Error is set.
	Result *BestQuoteResult `json:"result,omitempty"`
}
