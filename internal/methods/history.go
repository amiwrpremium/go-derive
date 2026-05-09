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
// Optional `params`: `start_timestamp`, `end_timestamp`,
// `instrument_name`, `page`, `page_size`. The configured subaccount
// is threaded through automatically when not present in `params`.
//
// The response is paginated; the second return value carries the
// total counts.
func (a *API) GetFundingHistory(ctx context.Context, params map[string]any) ([]types.FundingPayment, types.Page, error) {
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
// Optional `params`: `start_timestamp`, `end_timestamp`. The
// endpoint returns a bare array — there's no pagination wrapper.
func (a *API) GetLiquidationHistory(ctx context.Context, params map[string]any) ([]types.LiquidationAuction, error) {
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
	var resp []types.LiquidationAuction
	if err := a.call(ctx, "private/get_liquidation_history", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetOptionSettlementHistory returns the configured subaccount's
// past option-settlement events. Private.
//
// Optional `params`: `start_timestamp`, `end_timestamp`. The
// response is wrapped in a `settlements` array and is not
// paginated.
func (a *API) GetOptionSettlementHistory(ctx context.Context, params map[string]any) ([]types.OptionSettlement, error) {
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
	var resp struct {
		Settlements []types.OptionSettlement `json:"settlements"`
	}
	if err := a.call(ctx, "private/get_option_settlement_history", params, &resp); err != nil {
		return nil, err
	}
	return resp.Settlements, nil
}

// GetPublicOptionSettlementHistory returns the network-wide option
// settlement history. Public — no signer required.
//
// The endpoint paginates. The second return value carries the
// totals.
func (a *API) GetPublicOptionSettlementHistory(ctx context.Context, params map[string]any) ([]types.OptionSettlement, types.Page, error) {
	if params == nil {
		params = map[string]any{}
	}
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
// Required `params`: `period`, `start_timestamp`, `end_timestamp`.
// The response is wrapped in `subaccount_value_history` and is not
// paginated.
func (a *API) GetSubaccountValueHistory(ctx context.Context, params map[string]any) ([]types.SubaccountValueRecord, error) {
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
	var resp struct {
		SubaccountID           int64                         `json:"subaccount_id"`
		SubaccountValueHistory []types.SubaccountValueRecord `json:"subaccount_value_history"`
	}
	if err := a.call(ctx, "private/get_subaccount_value_history", params, &resp); err != nil {
		return nil, err
	}
	return resp.SubaccountValueHistory, nil
}

// GetERC20TransferHistory returns deposit / withdrawal-style ERC-20
// transfers attributed to the configured subaccount. Private.
//
// Optional `params`: `start_timestamp`, `end_timestamp`. The
// response is wrapped in `events` and is not paginated.
func (a *API) GetERC20TransferHistory(ctx context.Context, params map[string]any) ([]types.ERC20Transfer, error) {
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
// Optional `params`: `start_timestamp`, `end_timestamp`. The
// response is wrapped in `events` and is not paginated.
func (a *API) GetInterestHistory(ctx context.Context, params map[string]any) ([]types.InterestPayment, error) {
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
// Required `params`: `start_timestamp`, `end_timestamp`, `expiry`.
// The response carries pre-signed S3 URLs the caller can download
// directly to retrieve the archived records.
func (a *API) ExpiredAndCancelledHistory(ctx context.Context, params map[string]any) (*types.ExpiredAndCancelledExport, error) {
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
	var resp types.ExpiredAndCancelledExport
	if err := a.call(ctx, "private/expired_and_cancelled_history", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
