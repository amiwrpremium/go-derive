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

// SendQuote responds to an open RFQ with a maker quote. The
// payload carries the multi-leg priced quote plus the caller's
// pre-computed EIP-712 signature material. Private.
//
// The SDK does not yet sign quote payloads; the caller is
// responsible for populating [types.SendQuoteInput.Signature],
// [types.SendQuoteInput.Signer],
// [types.SendQuoteInput.SignatureExpirySec] and
// [types.SendQuoteInput.Nonce] before calling.
func (a *API) SendQuote(ctx context.Context, in types.SendQuoteInput) (types.Quote, error) {
	if err := a.requireSigner(); err != nil {
		return types.Quote{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.Quote{}, err
	}
	params := sendQuoteParams(in, a.Subaccount)
	var resp types.Quote
	if err := a.call(ctx, "private/send_quote", params, &resp); err != nil {
		return types.Quote{}, err
	}
	return resp, nil
}

// ExecuteQuote picks one quote response and trades against it. Used
// by the taker once `send_rfq` has surfaced acceptable quotes.
// Private.
//
// The SDK does not yet sign quote-execute payloads; the caller is
// responsible for the signature fields on [types.ExecuteQuoteInput].
// The response wraps a [types.Quote] and adds `rfq_filled_pct`.
func (a *API) ExecuteQuote(ctx context.Context, in types.ExecuteQuoteInput) (types.ExecuteQuoteResult, error) {
	if err := a.requireSigner(); err != nil {
		return types.ExecuteQuoteResult{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.ExecuteQuoteResult{}, err
	}
	params := executeQuoteParams(in, a.Subaccount)
	var resp types.ExecuteQuoteResult
	if err := a.call(ctx, "private/execute_quote", params, &resp); err != nil {
		return types.ExecuteQuoteResult{}, err
	}
	return resp, nil
}

// CancelQuote cancels one outstanding maker quote by id. Private.
func (a *API) CancelQuote(ctx context.Context, quoteID string) (types.Quote, error) {
	if err := a.requireSigner(); err != nil {
		return types.Quote{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.Quote{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"quote_id":      quoteID,
	}
	var resp types.Quote
	if err := a.call(ctx, "private/cancel_quote", params, &resp); err != nil {
		return types.Quote{}, err
	}
	return resp, nil
}

// ReplaceQuote cancels one outstanding maker quote and submits a
// replacement in a single round trip — the maker counterpart to
// [API.Replace] for orders. Private.
//
// The replacement carries the same signed-quote shape as
// [API.SendQuote]; the caller pre-signs and populates the
// signature fields on [types.SendQuoteInput] embedded in
// [types.ReplaceQuoteInput]. Exactly one of
// [types.ReplaceQuoteInput.QuoteIDToCancel] or
// [types.ReplaceQuoteInput.NonceToCancel] identifies the quote being
// replaced.
//
// The response carries the cancelled quote, the (optional)
// replacement quote, and the engine's error if the replacement was
// rejected.
func (a *API) ReplaceQuote(ctx context.Context, in types.ReplaceQuoteInput) (types.ReplaceQuoteResult, error) {
	if err := a.requireSigner(); err != nil {
		return types.ReplaceQuoteResult{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.ReplaceQuoteResult{}, err
	}
	params := sendQuoteParams(in.SendQuoteInput, a.Subaccount)
	if in.QuoteIDToCancel != "" {
		params["quote_id_to_cancel"] = in.QuoteIDToCancel
	}
	if in.NonceToCancel != 0 {
		params["nonce_to_cancel"] = in.NonceToCancel
	}
	var resp types.ReplaceQuoteResult
	if err := a.call(ctx, "private/replace_quote", params, &resp); err != nil {
		return types.ReplaceQuoteResult{}, err
	}
	return resp, nil
}

// CancelBatchQuotes cancels every quote matching the supplied
// filters. Private.
//
// All filters AND together. Returns the list of cancelled ids.
func (a *API) CancelBatchQuotes(ctx context.Context, filter types.CancelBatchInput) (types.CancelBatchResult, error) {
	if err := a.requireSigner(); err != nil {
		return types.CancelBatchResult{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.CancelBatchResult{}, err
	}
	params := cancelBatchParams(filter, a.Subaccount, true)
	var resp types.CancelBatchResult
	if err := a.call(ctx, "private/cancel_batch_quotes", params, &resp); err != nil {
		return types.CancelBatchResult{}, err
	}
	return resp, nil
}

// CancelBatchRFQs cancels every RFQ matching the supplied filters.
// Private.
//
// All filters AND together. The QuoteID field on
// [types.CancelBatchInput] is ignored by this endpoint. Returns the
// list of cancelled ids.
func (a *API) CancelBatchRFQs(ctx context.Context, filter types.CancelBatchInput) (types.CancelBatchResult, error) {
	if err := a.requireSigner(); err != nil {
		return types.CancelBatchResult{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.CancelBatchResult{}, err
	}
	params := cancelBatchParams(filter, a.Subaccount, false)
	var resp types.CancelBatchResult
	if err := a.call(ctx, "private/cancel_batch_rfqs", params, &resp); err != nil {
		return types.CancelBatchResult{}, err
	}
	return resp, nil
}

// RFQGetBestQuote returns the best quote the engine can match
// against an RFQ shape, plus margin-impact estimates for the
// would-be trade. Private.
//
// No signing required — the call is a pure lookup that the engine
// uses to surface the best executable price across whitelisted
// makers.
func (a *API) RFQGetBestQuote(ctx context.Context, in types.BestQuoteInput) (types.BestQuoteResult, error) {
	if err := a.requireSigner(); err != nil {
		return types.BestQuoteResult{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.BestQuoteResult{}, err
	}
	params := bestQuoteParams(in, a.Subaccount)
	var resp types.BestQuoteResult
	if err := a.call(ctx, "private/rfq_get_best_quote", params, &resp); err != nil {
		return types.BestQuoteResult{}, err
	}
	return resp, nil
}

func sendQuoteParams(in types.SendQuoteInput, defaultSubaccount int64) map[string]any {
	sub := in.SubaccountID
	if sub == 0 {
		sub = defaultSubaccount
	}
	params := map[string]any{
		"subaccount_id":        sub,
		"rfq_id":               in.RFQID,
		"direction":            in.Direction,
		"legs":                 in.Legs,
		"max_fee":              in.MaxFee,
		"nonce":                in.Nonce,
		"signature":            in.Signature,
		"signer":               in.Signer,
		"signature_expiry_sec": in.SignatureExpirySec,
	}
	if in.Label != "" {
		params["label"] = in.Label
	}
	if in.MMP {
		params["mmp"] = true
	}
	if in.Client != "" {
		params["client"] = in.Client
	}
	return params
}

func executeQuoteParams(in types.ExecuteQuoteInput, defaultSubaccount int64) map[string]any {
	sub := in.SubaccountID
	if sub == 0 {
		sub = defaultSubaccount
	}
	params := map[string]any{
		"subaccount_id":        sub,
		"rfq_id":               in.RFQID,
		"quote_id":             in.QuoteID,
		"direction":            in.Direction,
		"legs":                 in.Legs,
		"max_fee":              in.MaxFee,
		"nonce":                in.Nonce,
		"signature":            in.Signature,
		"signer":               in.Signer,
		"signature_expiry_sec": in.SignatureExpirySec,
	}
	if in.Label != "" {
		params["label"] = in.Label
	}
	if in.EnableTakerProtection {
		params["enable_taker_protection"] = true
	}
	if in.Client != "" {
		params["client"] = in.Client
	}
	return params
}

func bestQuoteParams(in types.BestQuoteInput, defaultSubaccount int64) map[string]any {
	sub := in.SubaccountID
	if sub == 0 {
		sub = defaultSubaccount
	}
	params := map[string]any{
		"subaccount_id": sub,
		"legs":          in.Legs,
	}
	if in.Direction != "" {
		params["direction"] = in.Direction
	}
	if in.PreferredDirection != "" {
		params["preferred_direction"] = in.PreferredDirection
	}
	if len(in.Counterparties) > 0 {
		params["counterparties"] = in.Counterparties
	}
	if in.Label != "" {
		params["label"] = in.Label
	}
	if in.Client != "" {
		params["client"] = in.Client
	}
	if !in.ExtraFee.IsZero() {
		params["extra_fee"] = in.ExtraFee
	}
	if !in.MaxTotalCost.IsZero() {
		params["max_total_cost"] = in.MaxTotalCost
	}
	if !in.MinTotalCost.IsZero() {
		params["min_total_cost"] = in.MinTotalCost
	}
	if !in.PartialFillStep.IsZero() {
		params["partial_fill_step"] = in.PartialFillStep
	}
	if in.ReferralCode != "" {
		params["referral_code"] = in.ReferralCode
	}
	if in.RFQID != "" {
		params["rfq_id"] = in.RFQID
	}
	return params
}

func cancelBatchParams(filter types.CancelBatchInput, defaultSubaccount int64, includeQuoteID bool) map[string]any {
	sub := filter.SubaccountID
	if sub == 0 {
		sub = defaultSubaccount
	}
	params := map[string]any{
		"subaccount_id": sub,
	}
	if filter.RFQID != "" {
		params["rfq_id"] = filter.RFQID
	}
	if includeQuoteID && filter.QuoteID != "" {
		params["quote_id"] = filter.QuoteID
	}
	if filter.Label != "" {
		params["label"] = filter.Label
	}
	if filter.Nonce != 0 {
		params["nonce"] = filter.Nonce
	}
	return params
}

// OrderQuotePublic is the unauthenticated counterpart of
// [API.OrderQuote] — it runs a hypothetical order through the
// matching engine without submitting and returns the engine's
// estimates for fill price, fee, and post-trade margin balance.
//
// Wraps `public/order_quote`. The endpoint still requires a fully
// signed order body (the connection is unauthenticated but the
// payload is not), so the SDK signs with the configured signer
// before sending. Same input and result shape as [API.OrderQuote].
func (a *API) OrderQuotePublic(ctx context.Context, in types.PlaceOrderInput) (types.OrderQuoteResult, error) {
	params, err := a.signedOrderParams(ctx, in)
	if err != nil {
		return types.OrderQuoteResult{}, err
	}
	var resp types.OrderQuoteResult
	if err := a.call(ctx, "public/order_quote", params, &resp); err != nil {
		return types.OrderQuoteResult{}, err
	}
	return resp, nil
}

// OrderQuote runs a hypothetical order through the matching engine
// without submitting and reports the engine's estimates for fill
// price, fee, and post-trade margin balance. Useful for
// pre-flighting orders against thin books where the user wants to
// know whether they'll clear margin before signing. Private.
//
// Input mirrors [types.PlaceOrderInput] — the shape `private/order`
// accepts. The SDK fills in subaccount id, nonce, signature, signer
// and expiry, exactly as for [API.PlaceOrder].
//
// The result shape mirrors `PrivateOrderQuoteResultSchema` in
// `derivexyz/cockpit/orderbook-types`.
func (a *API) OrderQuote(ctx context.Context, in types.PlaceOrderInput) (types.OrderQuoteResult, error) {
	params, err := a.signedOrderParams(ctx, in)
	if err != nil {
		return types.OrderQuoteResult{}, err
	}
	var resp types.OrderQuoteResult
	if err := a.call(ctx, "private/order_quote", params, &resp); err != nil {
		return types.OrderQuoteResult{}, err
	}
	return resp, nil
}
