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

// GetFundingRateHistory returns historical funding rate prints for one
// perpetual instrument over the requested window. Public.
//
// Required `params`: `instrument_name`. Optional: `start_timestamp`,
// `end_timestamp`, `period`.
func (a *API) GetFundingRateHistory(ctx context.Context, params map[string]any) ([]types.FundingRateHistoryItem, error) {
	var resp struct {
		FundingRateHistory []types.FundingRateHistoryItem `json:"funding_rate_history"`
	}
	if err := a.call(ctx, "public/get_funding_rate_history", params, &resp); err != nil {
		return nil, err
	}
	return resp.FundingRateHistory, nil
}

// GetSpotFeedHistory returns historical oracle spot prices for one
// currency over the requested window at the given period. Public.
//
// Required `params`: `currency`, `period`, `start_timestamp`,
// `end_timestamp`. Returns the currency the response is keyed against
// alongside the per-bucket samples.
func (a *API) GetSpotFeedHistory(ctx context.Context, params map[string]any) (currency string, items []types.SpotFeedHistoryItem, err error) {
	var resp struct {
		Currency        string                      `json:"currency"`
		SpotFeedHistory []types.SpotFeedHistoryItem `json:"spot_feed_history"`
	}
	if err := a.call(ctx, "public/get_spot_feed_history", params, &resp); err != nil {
		return "", nil, err
	}
	return resp.Currency, resp.SpotFeedHistory, nil
}

// GetLatestSignedFeeds returns the latest oracle signed-feed snapshot
// for every published currency, expiry, and feed type. Public.
//
// Optional `params`: `currency`. Pass nil to get every currency the
// venue publishes.
func (a *API) GetLatestSignedFeeds(ctx context.Context, params map[string]any) (types.SignedFeeds, error) {
	if params == nil {
		params = map[string]any{}
	}
	var resp types.SignedFeeds
	if err := a.call(ctx, "public/get_latest_signed_feeds", params, &resp); err != nil {
		return types.SignedFeeds{}, err
	}
	return resp, nil
}

// GetInterestRateHistory returns historical USDC borrow / supply
// APY prints over the requested window. Public.
//
// Required `params`: `from_timestamp_sec`, `to_timestamp_sec`.
// Optional: `page`, `page_size`. Paginated; the second return value
// carries the totals.
func (a *API) GetInterestRateHistory(ctx context.Context, params map[string]any) ([]types.InterestRateHistoryItem, types.Page, error) {
	var resp struct {
		InterestRates []types.InterestRateHistoryItem `json:"interest_rates"`
		Pagination    types.Page                      `json:"pagination"`
	}
	if err := a.call(ctx, "public/get_interest_rate_history", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.InterestRates, resp.Pagination, nil
}

// GetPerpImpactTWAP returns the time-weighted-average difference of
// mid price, ask-impact price, and bid-impact price versus spot for
// one currency's perpetual book over the requested window. Public.
//
// Required `params`: `currency`, `start_time`, `end_time`.
//
// The shape mirrors `PublicGetPerpImpactTwapResultSchema` in
// `derivexyz/cockpit/orderbook-types`.
func (a *API) GetPerpImpactTWAP(ctx context.Context, params map[string]any) (types.PerpImpactTWAP, error) {
	var resp types.PerpImpactTWAP
	if err := a.call(ctx, "public/get_perp_impact_twap", params, &resp); err != nil {
		return types.PerpImpactTWAP{}, err
	}
	return resp, nil
}
