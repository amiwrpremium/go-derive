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

	"github.com/amiwrpremium/go-derive"
)

// GetInstruments lists active instruments matching the filter. Public.
//
// Derive returns the result as a bare JSON array of instrument objects.
func (a *API) GetInstruments(ctx context.Context, currency string, kind derive.InstrumentType) ([]derive.Instrument, error) {
	params := map[string]any{}
	if currency != "" {
		params["currency"] = currency
	}
	if kind != "" {
		params["instrument_type"] = kind
	}
	params["expired"] = false
	var insts []derive.Instrument
	if err := a.call(ctx, "public/get_instruments", params, &insts); err != nil {
		return nil, err
	}
	return insts, nil
}

// GetInstrument fetches one instrument by name. Public.
func (a *API) GetInstrument(ctx context.Context, name string) (derive.Instrument, error) {
	var inst derive.Instrument
	err := a.call(ctx, "public/get_instrument", map[string]any{"instrument_name": name}, &inst)
	return inst, err
}

// GetTicker fetches the public ticker for one instrument. Public.
func (a *API) GetTicker(ctx context.Context, name string) (derive.Ticker, error) {
	var t derive.Ticker
	err := a.call(ctx, "public/get_ticker", map[string]any{"instrument_name": name}, &t)
	return t, err
}

// GetPublicTradeHistory returns recent trades on the instrument. Public.
func (a *API) GetPublicTradeHistory(ctx context.Context, instrument string, page derive.PageRequest) ([]derive.Trade, derive.Page, error) {
	params := map[string]any{"instrument_name": instrument}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Trades     []derive.Trade `json:"trades"`
		Pagination derive.Page    `json:"pagination"`
	}
	if err := a.call(ctx, "public/get_trade_history", params, &resp); err != nil {
		return nil, derive.Page{}, err
	}
	return resp.Trades, resp.Pagination, nil
}

// GetTime returns the server clock in milliseconds. Public.
func (a *API) GetTime(ctx context.Context) (int64, error) {
	var t int64
	err := a.call(ctx, "public/get_time", map[string]any{}, &t)
	return t, err
}

// GetCurrencies returns the list of supported currency names. Public.
//
// Derive's `public/get_all_currencies` result is a bare JSON array of
// rich currency objects (margin parameters, manager addresses, etc.);
// this method extracts the `currency` name field from each. Callers
// that need the full object should call the raw transport directly.
func (a *API) GetCurrencies(ctx context.Context) ([]string, error) {
	var raw []struct {
		Currency string `json:"currency"`
	}
	if err := a.call(ctx, "public/get_all_currencies", map[string]any{}, &raw); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(raw))
	for _, c := range raw {
		out = append(out, c.Currency)
	}
	return out, nil
}
