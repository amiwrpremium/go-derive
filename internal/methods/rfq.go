// Package methods is the shared implementation of every JSON-RPC method
// Derive exposes. Both pkg/rest.Client and pkg/ws.Client embed *API so that
// each method is defined exactly once, parameterised by the underlying
// transport.
//
// Public methods are unauthenticated; private methods require Signer to be
// non-nil. Private methods that mutate orders also use the Domain to sign
// the per-action EIP-712 hash.
package methods

import (
	"context"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// SendRFQ broadcasts a request-for-quote to market makers. Private.
func (a *API) SendRFQ(ctx context.Context, legs []types.RFQLeg, maxFee types.Decimal) (types.RFQ, error) {
	if err := a.requireSubaccount(); err != nil {
		return types.RFQ{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"legs":          legs,
		"max_total_fee": maxFee,
	}
	var rfq types.RFQ
	err := a.call(ctx, "private/send_rfq", params, &rfq)
	return rfq, err
}

// PollRFQs returns the status of recent RFQs initiated by this subaccount. Private.
func (a *API) PollRFQs(ctx context.Context) ([]types.RFQ, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var resp struct {
		RFQs []types.RFQ `json:"rfqs"`
	}
	err := a.call(ctx, "private/poll_rfqs", map[string]any{
		"subaccount_id": a.Subaccount,
	}, &resp)
	return resp.RFQs, err
}

// CancelRFQ cancels an outstanding RFQ. Private.
func (a *API) CancelRFQ(ctx context.Context, rfqID string) error {
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	return a.call(ctx, "private/cancel_rfq", map[string]any{
		"subaccount_id": a.Subaccount,
		"rfq_id":        rfqID,
	}, nil)
}

// RFQ flow on Derive: a taker issues `send_rfq`; makers see it on
// their `wallet.{addr}.rfqs` subscription and respond with
// `send_quote`; the taker picks one with `execute_quote`. The methods
// below cover every step of that flow plus the read / batch-cancel
// helpers.
//
// All wrappers thread `subaccount_id` automatically when not present
// in `params`. Sign payloads with the appropriate
// `pkg/auth.SignAction` (or `SignQuote`) before calling — the SDK
// passes `params` through verbatim.

// GetRFQs returns RFQs that match the supplied filters (open / done /
// status / time window / pagination). Private.
//
// Optional `params`: `rfq_id`, `status`, `from_timestamp`,
// `to_timestamp`, `page`, `page_size`. The configured subaccount is
// threaded through automatically.
func (a *API) GetRFQs(ctx context.Context, params map[string]any) ([]types.RFQ, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp struct {
		RFQs       []types.RFQ `json:"rfqs"`
		Pagination types.Page  `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_rfqs", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.RFQs, resp.Pagination, nil
}

// GetQuotes returns the maker-side full view of quotes the
// configured subaccount has issued or received. Private.
//
// Optional `params`: `rfq_id`, `quote_id`, `status`, `from_ts`,
// `to_ts`, `page`, `page_size`.
func (a *API) GetQuotes(ctx context.Context, params map[string]any) ([]types.Quote, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp struct {
		Quotes     []types.Quote `json:"quotes"`
		Pagination types.Page    `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_quotes", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Quotes, resp.Pagination, nil
}

// PollQuotes is the long-poll variant of [GetQuotes] — returns the
// taker-public view of quotes, suitable for makers / takers without
// access to the full signer-side body. Private.
//
// Same `params` shape as [GetQuotes]; only the response shape differs
// (`QuotePublic` instead of `Quote`).
func (a *API) PollQuotes(ctx context.Context, params map[string]any) ([]types.QuotePublic, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp struct {
		Quotes     []types.QuotePublic `json:"quotes"`
		Pagination types.Page          `json:"pagination"`
	}
	if err := a.call(ctx, "private/poll_quotes", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Quotes, resp.Pagination, nil
}

// SendQuote responds to an open RFQ with a maker quote. The signed
// payload covers the multi-leg quote price and a per-leg side
// direction. Private.
//
// Required `params` include `rfq_id`, the priced `legs`,
// `direction`, the signing fields (`signature`, `signer`,
// `signature_expiry_sec`, `nonce`), and `max_fee`.
func (a *API) SendQuote(ctx context.Context, params map[string]any) (*types.Quote, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp types.Quote
	if err := a.call(ctx, "private/send_quote", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExecuteQuote picks one quote response and trades against it. Used
// by the taker once `send_rfq` has surfaced acceptable quotes.
// Private.
//
// The response wraps a [types.Quote] and adds `rfq_filled_pct`.
func (a *API) ExecuteQuote(ctx context.Context, params map[string]any) (*types.ExecuteQuoteResult, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp types.ExecuteQuoteResult
	if err := a.call(ctx, "private/execute_quote", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelQuote cancels one outstanding maker quote by id. Private.
func (a *API) CancelQuote(ctx context.Context, params map[string]any) (*types.Quote, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp types.Quote
	if err := a.call(ctx, "private/cancel_quote", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReplaceQuote cancels one outstanding maker quote and submits a
// replacement in a single round trip — the maker counterpart to
// [API.Replace] for orders. Private.
//
// `params` should include `quote_id_to_cancel` (or
// `nonce_to_cancel`) plus the same fields [API.SendQuote] takes for
// the replacement (`rfq_id`, `direction`, `legs`, `max_fee`, the
// signing fields). The full param shape is documented at
// docs.derive.xyz.
//
// The response carries the cancelled quote, the (optional) replacement
// quote, and the engine's error if the replacement was rejected.
func (a *API) ReplaceQuote(ctx context.Context, params map[string]any) (*types.ReplaceQuoteResult, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp types.ReplaceQuoteResult
	if err := a.call(ctx, "private/replace_quote", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelBatchQuotes cancels every quote matching the supplied
// filters. Private.
//
// Optional `params`: `rfq_id`, `quote_id`, `label`, `nonce`. Returns
// the list of cancelled ids.
func (a *API) CancelBatchQuotes(ctx context.Context, params map[string]any) (*types.CancelBatchResult, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp types.CancelBatchResult
	if err := a.call(ctx, "private/cancel_batch_quotes", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelBatchRFQs cancels every RFQ matching the supplied filters.
// Private.
//
// Optional `params`: `rfq_id`, `label`, `nonce`. Returns the list of
// cancelled ids.
func (a *API) CancelBatchRFQs(ctx context.Context, params map[string]any) (*types.CancelBatchResult, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp types.CancelBatchResult
	if err := a.call(ctx, "private/cancel_batch_rfqs", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RFQGetBestQuote returns the best quote the engine can match
// against an RFQ shape, plus margin-impact estimates for the would-be
// trade. Private.
//
// Required `params`: `legs`, `direction`. Optional: `counterparties`,
// `label`, `client`, `extra_fee`, `max_total_cost`, `min_total_cost`,
// `partial_fill_step`, `preferred_direction`, `referral_code`,
// `rfq_id`.
func (a *API) RFQGetBestQuote(ctx context.Context, params map[string]any) (*types.BestQuoteResult, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp types.BestQuoteResult
	if err := a.call(ctx, "private/rfq_get_best_quote", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// OrderQuotePublic is the unauthenticated counterpart of
// [API.OrderQuote] — it runs a hypothetical order through the
// matching engine without submitting and returns the engine's
// estimates for fill price, fee, and post-trade margin balance.
//
// Wraps `public/order_quote`. Same `params` shape as
// [API.OrderQuote] (plus the same signing fields, which the engine
// uses for the simulation but does not validate); same response
// type. No signer or subaccount required — useful for pre-flight
// checks before signing anything.
func (a *API) OrderQuotePublic(ctx context.Context, params map[string]any) (*types.OrderQuoteResult, error) {
	if params == nil {
		params = map[string]any{}
	}
	var resp types.OrderQuoteResult
	if err := a.call(ctx, "public/order_quote", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// OrderQuote runs a hypothetical order through the matching engine
// without submitting and reports the engine's estimates for fill
// price, fee, and post-trade margin balance. Useful for
// pre-flighting orders against thin books where the user wants to
// know whether they'll clear margin before signing. Private.
//
// `params` mirror the shape `private/order` accepts.
//
// The shape mirrors `PrivateOrderQuoteResultSchema` in
// `derivexyz/cockpit/orderbook-types`.
func (a *API) OrderQuote(ctx context.Context, params map[string]any) (*types.OrderQuoteResult, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var resp types.OrderQuoteResult
	if err := a.call(ctx, "private/order_quote", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
