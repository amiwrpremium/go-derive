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
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// SetRFQModule is called by the client constructors to thread through
// the per-network RFQ module contract address.
func (a *API) SetRFQModule(addr common.Address) { a.rfqModule = addr }

// SendRFQ broadcasts a request-for-quote to market makers. Private.
//
// All filter / hint fields on [types.SendRFQInput] are optional; only
// [types.SendRFQInput.Legs] is required. Wire keys mirror
// docs.derive.xyz/reference/post_private-send-rfq.
func (a *API) SendRFQ(ctx context.Context, in types.SendRFQInput) (types.RFQ, error) {
	if err := a.requireSubaccount(); err != nil {
		return types.RFQ{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"legs":          in.Legs,
	}
	if len(in.Counterparties) > 0 {
		params["counterparties"] = in.Counterparties
	}
	if in.PreferredDirection != "" {
		params["preferred_direction"] = in.PreferredDirection
	}
	if in.ReducingDirection != "" {
		params["reducing_direction"] = in.ReducingDirection
	}
	if in.Label != "" {
		params["label"] = in.Label
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
	if in.Client != "" {
		params["client"] = in.Client
	}
	if in.ReferralCode != "" {
		params["referral_code"] = in.ReferralCode
	}
	if !in.ExtraFee.IsZero() {
		params["extra_fee"] = in.ExtraFee
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
func (a *API) CancelRFQ(ctx context.Context, in types.CancelRFQInput) error {
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	return a.call(ctx, "private/cancel_rfq", map[string]any{
		"subaccount_id": a.Subaccount,
		"rfq_id":        in.RFQID,
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
// The configured subaccount is threaded through automatically. The
// query's `[FromTimestamp, ToTimestamp]` window filters on each
// RFQ's `last_update_timestamp`.
func (a *API) GetRFQs(ctx context.Context, q types.RFQsQuery, page types.PageRequest) ([]types.RFQ, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
	}
	if q.RFQID != "" {
		params["rfq_id"] = q.RFQID
	}
	if q.Status != "" {
		params["status"] = q.Status
	}
	addRFQWindow(params, q.HistoryWindow)
	addPaging(params, page)
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
func (a *API) GetQuotes(ctx context.Context, q types.QuotesQuery, page types.PageRequest) ([]types.Quote, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	params := quotesQueryParams(q, a.Subaccount)
	addPaging(params, page)
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
// Same query shape as [GetQuotes]; only the response shape differs
// (`QuotePublic` instead of `Quote`).
func (a *API) PollQuotes(ctx context.Context, q types.PollQuotesQuery, page types.PageRequest) ([]types.QuotePublic, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	params := quotesQueryParams(q, a.Subaccount)
	addPaging(params, page)
	var resp struct {
		Quotes     []types.QuotePublic `json:"quotes"`
		Pagination types.Page          `json:"pagination"`
	}
	if err := a.call(ctx, "private/poll_quotes", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Quotes, resp.Pagination, nil
}

func addRFQWindow(params map[string]any, w types.HistoryWindow) {
	if !w.StartTimestamp.Time().IsZero() {
		params["from_timestamp"] = w.StartTimestamp.Millis()
	}
	if !w.EndTimestamp.Time().IsZero() {
		params["to_timestamp"] = w.EndTimestamp.Millis()
	}
}

func quotesQueryParams(q types.QuotesQuery, defaultSubaccount int64) map[string]any {
	params := map[string]any{
		"subaccount_id": defaultSubaccount,
	}
	if q.RFQID != "" {
		params["rfq_id"] = q.RFQID
	}
	if q.QuoteID != "" {
		params["quote_id"] = q.QuoteID
	}
	if q.Status != "" {
		params["status"] = q.Status
	}
	addRFQWindow(params, q.HistoryWindow)
	return params
}

// SendQuote responds to an open RFQ with a maker quote. The SDK
// signs the per-quote EIP-712 payload internally; the caller
// supplies only the business fields on [types.SendQuoteInput]
// (RFQ id, direction, legs, max fee, optional metadata). Private.
//
// Each leg must carry both the engine-facing fields and the
// on-chain identifiers `Asset` + `SubID` (used by the RFQ module's
// per-leg hash). Retrieve them via `public/get_instrument` once
// per instrument at startup.
func (a *API) SendQuote(ctx context.Context, in types.SendQuoteInput) (types.Quote, error) {
	params, err := a.signedQuoteParams(ctx, in)
	if err != nil {
		return types.Quote{}, err
	}
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
// The SDK signs the per-execute EIP-712 payload internally,
// inverting the global direction when computing the per-leg signed
// amount (the taker takes the opposite side of the maker quote).
// The response wraps a [types.Quote] and adds `rfq_filled_pct`.
func (a *API) ExecuteQuote(ctx context.Context, in types.ExecuteQuoteInput) (types.ExecuteQuoteResult, error) {
	params, err := a.signedExecuteQuoteParams(ctx, in)
	if err != nil {
		return types.ExecuteQuoteResult{}, err
	}
	var resp types.ExecuteQuoteResult
	if err := a.call(ctx, "private/execute_quote", params, &resp); err != nil {
		return types.ExecuteQuoteResult{}, err
	}
	return resp, nil
}

// CancelQuote cancels one outstanding maker quote by id. Private.
func (a *API) CancelQuote(ctx context.Context, in types.CancelQuoteInput) (types.Quote, error) {
	if err := a.requireSigner(); err != nil {
		return types.Quote{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.Quote{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"quote_id":      in.QuoteID,
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
// The SDK signs the replacement quote exactly like
// [API.SendQuote]. Exactly one of
// [types.ReplaceQuoteInput.QuoteIDToCancel] or
// [types.ReplaceQuoteInput.NonceToCancel] identifies the quote
// being replaced.
//
// The response carries the cancelled quote, the (optional)
// replacement quote, and the engine's error if the replacement was
// rejected.
func (a *API) ReplaceQuote(ctx context.Context, in types.ReplaceQuoteInput) (types.ReplaceQuoteResult, error) {
	params, err := a.signedQuoteParams(ctx, in.SendQuoteInput)
	if err != nil {
		return types.ReplaceQuoteResult{}, err
	}
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

// signedQuoteParams is the shared signing block for the
// quote-submission endpoints (private/send_quote,
// private/replace_quote). It hashes the per-quote RFQ module
// payload, wraps it in the EIP-712 ActionData envelope, signs,
// and returns the wire params map.
func (a *API) signedQuoteParams(ctx context.Context, in types.SendQuoteInput) (map[string]any, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}

	if err := a.resolveQuoteLegs(ctx, in.Legs); err != nil {
		return nil, err
	}

	sub := in.SubaccountID
	if sub == 0 {
		sub = a.Subaccount
	}
	nonce := a.Nonces.Next()
	expiry := time.Now().Unix() + a.SignatureExpiry

	legs := convertLegs(in.Legs)
	module := auth.RFQQuoteModuleData{
		GlobalDirection: in.Direction,
		MaxFee:          in.MaxFee.Inner(),
		Legs:            legs,
	}
	dataHash, err := module.Hash()
	if err != nil {
		return nil, err
	}

	action := auth.ActionData{
		SubaccountID: sub,
		Nonce:        nonce,
		Module:       a.rfqModule,
		Data:         dataHash,
		Expiry:       expiry,
		Owner:        a.Signer.OwnerAddress(),
		Signer:       a.Signer.SessionAddress(),
	}
	sig, err := a.Signer.SignAction(ctx, a.Domain, action)
	if err != nil {
		return nil, err
	}

	params := map[string]any{
		"subaccount_id":        sub,
		"rfq_id":               in.RFQID,
		"direction":            in.Direction,
		"legs":                 in.Legs,
		"max_fee":              in.MaxFee,
		"nonce":                nonce,
		"signature":            sig.Hex(),
		"signer":               a.Signer.SessionAddress().Hex(),
		"signature_expiry_sec": expiry,
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
	return params, nil
}

// signedExecuteQuoteParams is the shared signing block for
// private/execute_quote. The leg encoding inverts the global
// direction (the taker takes the opposite side of the maker quote).
func (a *API) signedExecuteQuoteParams(ctx context.Context, in types.ExecuteQuoteInput) (map[string]any, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}

	if err := a.resolveQuoteLegs(ctx, in.Legs); err != nil {
		return nil, err
	}

	sub := in.SubaccountID
	if sub == 0 {
		sub = a.Subaccount
	}
	nonce := a.Nonces.Next()
	expiry := time.Now().Unix() + a.SignatureExpiry

	legs := convertLegs(in.Legs)
	module := auth.RFQExecuteModuleData{
		GlobalDirection: in.Direction,
		MaxFee:          in.MaxFee.Inner(),
		Legs:            legs,
	}
	dataHash, err := module.Hash()
	if err != nil {
		return nil, err
	}

	action := auth.ActionData{
		SubaccountID: sub,
		Nonce:        nonce,
		Module:       a.rfqModule,
		Data:         dataHash,
		Expiry:       expiry,
		Owner:        a.Signer.OwnerAddress(),
		Signer:       a.Signer.SessionAddress(),
	}
	sig, err := a.Signer.SignAction(ctx, a.Domain, action)
	if err != nil {
		return nil, err
	}

	params := map[string]any{
		"subaccount_id":        sub,
		"rfq_id":               in.RFQID,
		"quote_id":             in.QuoteID,
		"direction":            in.Direction,
		"legs":                 in.Legs,
		"max_fee":              in.MaxFee,
		"nonce":                nonce,
		"signature":            sig.Hex(),
		"signer":               a.Signer.SessionAddress().Hex(),
		"signature_expiry_sec": expiry,
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
	return params, nil
}

// convertLegs maps the public types.QuoteLeg shape (used on the
// wire and on inputs) to the auth.RFQQuoteLeg shape (used inside
// the signing primitives). The wire fields and on-chain
// identifiers travel together on a single QuoteLeg, so the
// translation is a 1:1 field copy.
// resolveQuoteLegs fills in Asset/SubID on any leg the caller left
// zero, using the instrument cache (or fetching public/get_instrument
// on miss). Mutates the supplied slice in place — legs are
// reference-shared with the wire params map but Asset/SubID are
// json:"-" so the mutation is invisible on the wire.
func (a *API) resolveQuoteLegs(ctx context.Context, legs []types.QuoteLeg) error {
	for i := range legs {
		if !legs[i].Asset.IsZero() || legs[i].InstrumentName == "" {
			continue
		}
		meta, err := a.resolveInstrument(ctx, legs[i].InstrumentName)
		if err != nil {
			return err
		}
		legs[i].Asset = meta.Asset
		legs[i].SubID = meta.SubID
	}
	return nil
}

func convertLegs(legs []types.QuoteLeg) []auth.RFQQuoteLeg {
	out := make([]auth.RFQQuoteLeg, len(legs))
	for i := range legs {
		out[i] = auth.RFQQuoteLeg{
			Asset:     common.Address(legs[i].Asset),
			SubID:     legs[i].SubID,
			Direction: legs[i].Direction,
			Amount:    legs[i].Amount.Inner(),
			Price:     legs[i].Price.Inner(),
		}
	}
	return out
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
