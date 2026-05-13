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

// GetFundingHistory returns funding payments received / paid by the
// configured subaccount over the requested window. Private.
//
// The response is paginated; the second return value carries the
// total counts.
func (a *API) GetFundingHistory(ctx context.Context, q types.FundingHistoryQuery, page types.PageRequest) ([]types.FundingPayment, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
	}
	if q.InstrumentName != "" {
		params["instrument_name"] = q.InstrumentName
	}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	}
	addHistoryWindow(params, q.HistoryWindow)
	addPaging(params, page)
	var resp struct {
		Events     []types.FundingPayment `json:"events"`
		Pagination types.Page             `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_funding_history", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Events, resp.Pagination, nil
}

// GetLiquidationHistory returns the configured subaccount's past
// liquidation events. Private.
//
// The endpoint returns a bare array — there's no pagination wrapper.
// When [types.LiquidationHistoryQuery.Wallet] is set the engine
// queries across all of its subaccounts and ignores
// `subaccount_id`.
func (a *API) GetLiquidationHistory(ctx context.Context, q types.LiquidationHistoryQuery) ([]types.LiquidationAuction, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
	}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	}
	addHistoryWindow(params, q.HistoryWindow)
	var resp []types.LiquidationAuction
	if err := a.call(ctx, "private/get_liquidation_history", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetLiquidatorHistory returns auction bids placed by the configured
// subaccount as a liquidator. Private.
//
// The response is paginated; the second return value carries the
// totals.
func (a *API) GetLiquidatorHistory(ctx context.Context, q types.LiquidatorHistoryQuery, page types.PageRequest) ([]types.AuctionBid, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
	}
	addHistoryWindow(params, q.HistoryWindow)
	addPaging(params, page)
	var resp struct {
		Bids       []types.AuctionBid `json:"bids"`
		Pagination types.Page         `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_liquidator_history", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Bids, resp.Pagination, nil
}

// GetOptionSettlementHistory returns the configured subaccount's
// past option-settlement events. Private.
//
// The endpoint takes no time window — query by account identity
// only. The response is wrapped in a `settlements` array and is not
// paginated.
func (a *API) GetOptionSettlementHistory(ctx context.Context, q types.OptionSettlementHistoryQuery) ([]types.OptionSettlement, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
	}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	}
	var resp struct {
		Settlements []types.OptionSettlement `json:"settlements"`
	}
	if err := a.call(ctx, "private/get_option_settlement_history", params, &resp); err != nil {
		return nil, err
	}
	return resp.Settlements, nil
}

// GetPublicLiquidationHistory returns the network-wide liquidation
// history. Public — no signer required.
//
// Counterpart to the private [API.GetLiquidationHistory] (configured
// subaccount only). The wire schemas differ in their wrapper field
// name — public uses `auctions`, private uses a bare array — but
// the per-record shape ([types.LiquidationAuction]) is the same.
func (a *API) GetPublicLiquidationHistory(ctx context.Context, q types.LiquidationHistoryQuery, page types.PageRequest) ([]types.LiquidationAuction, types.Page, error) {
	params := map[string]any{}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	}
	addHistoryWindow(params, q.HistoryWindow)
	addPaging(params, page)
	var resp struct {
		Auctions   []types.LiquidationAuction `json:"auctions"`
		Pagination types.Page                 `json:"pagination"`
	}
	if err := a.call(ctx, "public/get_liquidation_history", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Auctions, resp.Pagination, nil
}

// GetPublicOptionSettlementHistory returns the network-wide option
// settlement history. Public — no signer required.
//
// The endpoint takes no time window. The response is paginated;
// the second return value carries the totals.
func (a *API) GetPublicOptionSettlementHistory(ctx context.Context, q types.OptionSettlementHistoryQuery, page types.PageRequest) ([]types.OptionSettlement, types.Page, error) {
	params := map[string]any{}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	} else if q.SubaccountID != 0 {
		params["subaccount_id"] = q.SubaccountID
	}
	addPaging(params, page)
	var resp struct {
		Settlements []types.OptionSettlement `json:"settlements"`
		Pagination  types.Page               `json:"pagination"`
	}
	if err := a.call(ctx, "public/get_option_settlement_history", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Settlements, resp.Pagination, nil
}

// GetSubaccountValueHistory returns the equity-curve series for the
// configured subaccount. Private.
//
// Returns the subaccount id the engine echoed back alongside the
// per-bucket samples; the response is not paginated.
func (a *API) GetSubaccountValueHistory(ctx context.Context, q types.SubaccountValueHistoryQuery) (subaccountID int64, history []types.SubaccountValueRecord, err error) {
	if err := a.requireSigner(); err != nil {
		return 0, nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return 0, nil, err
	}
	params := map[string]any{
		"subaccount_id":   a.Subaccount,
		"period":          q.PeriodSec,
		"start_timestamp": q.StartTimestamp.Millis(),
		"end_timestamp":   q.EndTimestamp.Millis(),
	}
	var resp struct {
		SubaccountID           int64                         `json:"subaccount_id"`
		SubaccountValueHistory []types.SubaccountValueRecord `json:"subaccount_value_history"`
	}
	if err := a.call(ctx, "private/get_subaccount_value_history", params, &resp); err != nil {
		return 0, nil, err
	}
	return resp.SubaccountID, resp.SubaccountValueHistory, nil
}

// GetERC20TransferHistory returns deposit / withdrawal-style ERC-20
// transfers attributed to the configured subaccount. Private.
//
// The response is wrapped in `events` and is not paginated.
func (a *API) GetERC20TransferHistory(ctx context.Context, q types.ERC20TransferHistoryQuery) ([]types.ERC20Transfer, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
	}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	}
	addHistoryWindow(params, q.HistoryWindow)
	var resp struct {
		Events []types.ERC20Transfer `json:"events"`
	}
	if err := a.call(ctx, "private/get_erc20_transfer_history", params, &resp); err != nil {
		return nil, err
	}
	return resp.Events, nil
}

// GetInterestHistory returns the configured subaccount's interest
// charges and rebates. Private.
//
// The response is wrapped in `events` and is not paginated.
func (a *API) GetInterestHistory(ctx context.Context, q types.InterestHistoryQuery) ([]types.InterestPayment, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
	}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	}
	addHistoryWindow(params, q.HistoryWindow)
	var resp struct {
		Events []types.InterestPayment `json:"events"`
	}
	if err := a.call(ctx, "private/get_interest_history", params, &resp); err != nil {
		return nil, err
	}
	return resp.Events, nil
}

// ExpiredAndCancelledHistory triggers an archive export of the
// configured subaccount's expired and cancelled orders. Private.
//
// The response carries pre-signed S3 URLs the caller can download
// directly to retrieve the archived records.
func (a *API) ExpiredAndCancelledHistory(ctx context.Context, in types.ExpiredAndCancelledHistoryInput) (types.ExpiredAndCancelledExport, error) {
	if err := a.requireSigner(); err != nil {
		return types.ExpiredAndCancelledExport{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.ExpiredAndCancelledExport{}, err
	}
	sub := in.SubaccountID
	if sub == 0 {
		sub = a.Subaccount
	}
	params := map[string]any{
		"subaccount_id":   sub,
		"wallet":          in.Wallet,
		"start_timestamp": in.StartTimestamp.Millis(),
		"end_timestamp":   in.EndTimestamp.Millis(),
		"expiry":          in.ExpirySec,
	}
	var resp types.ExpiredAndCancelledExport
	if err := a.call(ctx, "private/expired_and_cancelled_history", params, &resp); err != nil {
		return types.ExpiredAndCancelledExport{}, err
	}
	return resp, nil
}

func addHistoryWindow(params map[string]any, w types.HistoryWindow) {
	if !w.StartTimestamp.Time().IsZero() {
		params["start_timestamp"] = w.StartTimestamp.Millis()
	}
	if !w.EndTimestamp.Time().IsZero() {
		params["end_timestamp"] = w.EndTimestamp.Millis()
	}
}

func addPaging(params map[string]any, p types.PageRequest) {
	if p.Page > 0 {
		params["page"] = p.Page
	}
	if p.PageSize > 0 {
		params["page_size"] = p.PageSize
	}
}
