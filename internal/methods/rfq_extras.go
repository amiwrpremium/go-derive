// Package methods.
package methods

import (
	"context"
	"encoding/json"
)

// RFQ flow on Derive: a taker issues `send_rfq`; makers see it on their
// `wallet.{addr}.rfqs` subscription and respond with `send_quote`; the
// taker picks one with `execute_quote`. The methods below cover every
// step of that flow plus the read / batch-cancel helpers.
//
// All wrappers are thin — parameter shapes follow Derive's documentation
// (multi-leg signed payloads), so the wrappers take `map[string]any` and
// return `json.RawMessage`. Sign the payload with `pkg/auth.SignAction`
// before calling.

// GetRFQs returns the configured subaccount's outstanding (open / done)
// RFQs.
func (a *API) GetRFQs(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	var raw json.RawMessage
	err := a.call(ctx, "private/get_rfqs", params, &raw)
	return raw, err
}

// GetQuotes returns quotes the configured subaccount has issued or
// received against open RFQs.
func (a *API) GetQuotes(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	var raw json.RawMessage
	err := a.call(ctx, "private/get_quotes", params, &raw)
	return raw, err
}

// PollQuotes is the long-poll variant of GetQuotes — used by makers who
// want to be woken on new RFQs without holding a WebSocket open.
func (a *API) PollQuotes(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	var raw json.RawMessage
	err := a.call(ctx, "private/poll_quotes", params, &raw)
	return raw, err
}

// SendQuote responds to an open RFQ with a maker quote. The signed
// payload covers the multi-leg quote price and a per-leg side direction.
//
// Required params include the RFQ id, the per-leg quote prices, and the
// signature/nonce/expiry triple. Private.
func (a *API) SendQuote(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	var raw json.RawMessage
	err := a.call(ctx, "private/send_quote", params, &raw)
	return raw, err
}

// ExecuteQuote picks one quote response and trades against it. Used by
// the taker once `send_rfq` has surfaced acceptable quotes.
func (a *API) ExecuteQuote(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	var raw json.RawMessage
	err := a.call(ctx, "private/execute_quote", params, &raw)
	return raw, err
}

// CancelQuote cancels one outstanding maker quote by id.
func (a *API) CancelQuote(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	var raw json.RawMessage
	err := a.call(ctx, "private/cancel_quote", params, &raw)
	return raw, err
}

// CancelBatchQuotes cancels every quote whose id appears in `quote_ids`,
// or every open quote on the subaccount when the field is omitted.
func (a *API) CancelBatchQuotes(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	var raw json.RawMessage
	err := a.call(ctx, "private/cancel_batch_quotes", params, &raw)
	return raw, err
}

// CancelBatchRFQs cancels every RFQ whose id appears in `rfq_ids`, or
// every open RFQ on the subaccount when the field is omitted.
func (a *API) CancelBatchRFQs(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	var raw json.RawMessage
	err := a.call(ctx, "private/cancel_batch_rfqs", params, &raw)
	return raw, err
}

// RFQGetBestQuote returns the best quote currently outstanding on one
// RFQ — the helper a taker uses to pick a counterparty before calling
// ExecuteQuote.
func (a *API) RFQGetBestQuote(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	var raw json.RawMessage
	err := a.call(ctx, "private/rfq_get_best_quote", params, &raw)
	return raw, err
}

// OrderQuote routes an order through the RFQ matching path instead of
// the central order book. Useful for instruments with thin books where
// makers respond on demand.
func (a *API) OrderQuote(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	var raw json.RawMessage
	err := a.call(ctx, "private/order_quote", params, &raw)
	return raw, err
}
