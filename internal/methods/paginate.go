// Package methods is the shared implementation of every JSON-RPC method
// Derive exposes.
//
// This file holds the "*All" companion methods that exhaust pagination on
// every paginated endpoint the SDK exposes. Each wrapper is a thin
// closure over [types.Paginate]; the underlying paginated method stays
// available unchanged for callers who want page-by-page control.
package methods

import (
	"context"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// GetOrdersAll exhausts pagination on [API.GetOrders].
func (a *API) GetOrdersAll(ctx context.Context, filter *types.GetOrdersFilter, opts types.PaginateOptions) ([]types.Order, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.Order, types.Page, error) {
		return a.GetOrders(ctx, page, filter)
	})
}

// GetOrderHistoryAll exhausts pagination on [API.GetOrderHistory].
func (a *API) GetOrderHistoryAll(ctx context.Context, q types.OrderHistoryQuery, opts types.PaginateOptions) ([]types.Order, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.Order, types.Page, error) {
		return a.GetOrderHistory(ctx, page, q)
	})
}

// GetTradeHistoryAll exhausts pagination on the private [API.GetTradeHistory].
func (a *API) GetTradeHistoryAll(ctx context.Context, q types.TradeHistoryQuery, opts types.PaginateOptions) ([]types.Trade, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.Trade, types.Page, error) {
		return a.GetTradeHistory(ctx, q, page)
	})
}

// GetPublicTradeHistoryAll exhausts pagination on [API.GetPublicTradeHistory].
func (a *API) GetPublicTradeHistoryAll(ctx context.Context, q types.PublicTradeHistoryQuery, opts types.PaginateOptions) ([]types.Trade, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.Trade, types.Page, error) {
		return a.GetPublicTradeHistory(ctx, q, page)
	})
}

// GetDepositHistoryAll exhausts pagination on [API.GetDepositHistory].
func (a *API) GetDepositHistoryAll(ctx context.Context, opts types.PaginateOptions) ([]types.DepositTx, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.DepositTx, types.Page, error) {
		return a.GetDepositHistory(ctx, page)
	})
}

// GetWithdrawalHistoryAll exhausts pagination on [API.GetWithdrawalHistory].
func (a *API) GetWithdrawalHistoryAll(ctx context.Context, opts types.PaginateOptions) ([]types.WithdrawTx, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.WithdrawTx, types.Page, error) {
		return a.GetWithdrawalHistory(ctx, page)
	})
}

// GetAllInstrumentsAll exhausts pagination on [API.GetAllInstruments].
//
// The `includeExpired` flag is threaded through to every fetch. Use it
// to opt into expired instruments for historical lookups; leave it
// false for trading-side cache warming.
func (a *API) GetAllInstrumentsAll(ctx context.Context, kind enums.InstrumentType, includeExpired bool, opts types.PaginateOptions) ([]types.Instrument, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.Instrument, types.Page, error) {
		return a.GetAllInstruments(ctx, kind, includeExpired, page)
	})
}

// GetPublicOptionSettlementHistoryAll exhausts pagination on
// [API.GetPublicOptionSettlementHistory].
func (a *API) GetPublicOptionSettlementHistoryAll(ctx context.Context, q types.OptionSettlementHistoryQuery, opts types.PaginateOptions) ([]types.OptionSettlement, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.OptionSettlement, types.Page, error) {
		return a.GetPublicOptionSettlementHistory(ctx, q, page)
	})
}

// GetPublicLiquidationHistoryAll exhausts pagination on
// [API.GetPublicLiquidationHistory].
func (a *API) GetPublicLiquidationHistoryAll(ctx context.Context, q types.LiquidationHistoryQuery, opts types.PaginateOptions) ([]types.LiquidationAuction, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.LiquidationAuction, types.Page, error) {
		return a.GetPublicLiquidationHistory(ctx, q, page)
	})
}

// GetLiquidatorHistoryAll exhausts pagination on [API.GetLiquidatorHistory].
func (a *API) GetLiquidatorHistoryAll(ctx context.Context, q types.LiquidatorHistoryQuery, opts types.PaginateOptions) ([]types.AuctionBid, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.AuctionBid, types.Page, error) {
		return a.GetLiquidatorHistory(ctx, q, page)
	})
}

// GetFundingHistoryAll exhausts pagination on [API.GetFundingHistory].
func (a *API) GetFundingHistoryAll(ctx context.Context, q types.FundingHistoryQuery, opts types.PaginateOptions) ([]types.FundingPayment, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.FundingPayment, types.Page, error) {
		return a.GetFundingHistory(ctx, q, page)
	})
}

// GetInterestRateHistoryAll exhausts pagination on [API.GetInterestRateHistory].
func (a *API) GetInterestRateHistoryAll(ctx context.Context, q types.InterestRateHistoryQuery, opts types.PaginateOptions) ([]types.InterestRateHistoryItem, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.InterestRateHistoryItem, types.Page, error) {
		return a.GetInterestRateHistory(ctx, q, page)
	})
}

// GetNotificationsAll exhausts pagination on [API.GetNotifications].
func (a *API) GetNotificationsAll(ctx context.Context, q types.NotificationsQuery, opts types.PaginateOptions) ([]types.Notification, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.Notification, types.Page, error) {
		return a.GetNotifications(ctx, q, page)
	})
}

// GetVaultShareAll exhausts pagination on [API.GetVaultShare].
func (a *API) GetVaultShareAll(ctx context.Context, q types.VaultShareQuery, opts types.PaginateOptions) ([]types.VaultShare, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.VaultShare, types.Page, error) {
		return a.GetVaultShare(ctx, q, page)
	})
}

// GetRFQsAll exhausts pagination on [API.GetRFQs].
func (a *API) GetRFQsAll(ctx context.Context, q types.RFQsQuery, opts types.PaginateOptions) ([]types.RFQ, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.RFQ, types.Page, error) {
		return a.GetRFQs(ctx, q, page)
	})
}

// GetQuotesAll exhausts pagination on [API.GetQuotes].
func (a *API) GetQuotesAll(ctx context.Context, q types.QuotesQuery, opts types.PaginateOptions) ([]types.Quote, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.Quote, types.Page, error) {
		return a.GetQuotes(ctx, q, page)
	})
}

// PollQuotesAll exhausts pagination on [API.PollQuotes].
func (a *API) PollQuotesAll(ctx context.Context, q types.PollQuotesQuery, opts types.PaginateOptions) ([]types.QuotePublic, error) {
	return types.Paginate(ctx, opts, func(ctx context.Context, page types.PageRequest) ([]types.QuotePublic, types.Page, error) {
		return a.PollQuotes(ctx, q, page)
	})
}
