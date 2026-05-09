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

// GetTradeHistory paginates the user's fills. Private.
func (a *API) GetTradeHistory(ctx context.Context, page derive.PageRequest) ([]derive.Trade, derive.Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, derive.Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
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
	if err := a.call(ctx, "private/get_trade_history", params, &resp); err != nil {
		return nil, derive.Page{}, err
	}
	return resp.Trades, resp.Pagination, nil
}

// GetDepositHistory paginates deposit transactions. Private.
func (a *API) GetDepositHistory(ctx context.Context, page derive.PageRequest) ([]derive.DepositTx, derive.Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, derive.Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Events     []derive.DepositTx `json:"events"`
		Pagination derive.Page        `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_deposit_history", params, &resp); err != nil {
		return nil, derive.Page{}, err
	}
	return resp.Events, resp.Pagination, nil
}

// GetWithdrawalHistory paginates withdrawal transactions. Private.
func (a *API) GetWithdrawalHistory(ctx context.Context, page derive.PageRequest) ([]derive.WithdrawTx, derive.Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, derive.Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Events     []derive.WithdrawTx `json:"events"`
		Pagination derive.Page         `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_withdrawal_history", params, &resp); err != nil {
		return nil, derive.Page{}, err
	}
	return resp.Events, resp.Pagination, nil
}
