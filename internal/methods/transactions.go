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

// GetTradeHistory paginates the user's fills. Private.
//
// All filters on the query are optional. Pass Wallet to span every
// subaccount under that wallet; when Wallet is empty the
// client-configured subaccount is used. InstrumentName, OrderID, and
// QuoteID narrow the result further; QuoteID accepts a concrete UUID
// or the engine's enum strings "is_quote" / "is_not_quote".
func (a *API) GetTradeHistory(ctx context.Context, q types.TradeHistoryQuery, page types.PageRequest) ([]types.Trade, types.Page, error) {
	params := map[string]any{}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	} else {
		if err := a.requireSubaccount(); err != nil {
			return nil, types.Page{}, err
		}
		params["subaccount_id"] = a.Subaccount
	}
	if q.InstrumentName != "" {
		params["instrument_name"] = q.InstrumentName
	}
	if q.OrderID != "" {
		params["order_id"] = q.OrderID
	}
	if q.QuoteID != "" {
		params["quote_id"] = q.QuoteID
	}
	if !q.FromTimestamp.Time().IsZero() {
		params["from_timestamp"] = q.FromTimestamp.Millis()
	}
	if !q.ToTimestamp.Time().IsZero() {
		params["to_timestamp"] = q.ToTimestamp.Millis()
	}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Trades     []types.Trade `json:"trades"`
		Pagination types.Page    `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_trade_history", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Trades, resp.Pagination, nil
}

// GetDepositHistory paginates deposit transactions. Private.
func (a *API) GetDepositHistory(ctx context.Context, page types.PageRequest) ([]types.DepositTx, types.Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Events     []types.DepositTx `json:"events"`
		Pagination types.Page        `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_deposit_history", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Events, resp.Pagination, nil
}

// GetWithdrawalHistory paginates withdrawal transactions. Private.
func (a *API) GetWithdrawalHistory(ctx context.Context, page types.PageRequest) ([]types.WithdrawTx, types.Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Events     []types.WithdrawTx `json:"events"`
		Pagination types.Page         `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_withdrawal_history", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Events, resp.Pagination, nil
}

// GetTransaction returns the on-chain status / receipt of one Derive
// transaction by its server-side id. Public.
func (a *API) GetTransaction(ctx context.Context, transactionID string) (*types.Transaction, error) {
	var resp types.Transaction
	if err := a.call(ctx, "public/get_transaction", map[string]any{"transaction_id": transactionID}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
