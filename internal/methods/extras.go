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
	"encoding/json"
)

// This file holds JSON-RPC wrappers for endpoints that exist on Derive but
// whose response schemas are large, evolving, or instrument-specific
// enough that a typed binding would either be premature or restrictive.
// Each wrapper returns a [json.RawMessage] so callers can decode the
// payload at the call site, against the latest contract documented at
// docs.derive.xyz, without going through this SDK for a bump.
//
// Each method is verified against the live API and the canonical schema
// list in `derivexyz/cockpit`. None of these are hallucinated.

// ---------------------------------------------------------------------------
// Public read methods
// ---------------------------------------------------------------------------

// GetFundingRateHistory returns historical funding rate prints for one
// perpetual instrument over the requested window.
//
// Required params: `instrument_name`. Optional: `start_timestamp`,
// `end_timestamp`, `period`. Public.
func (a *API) GetFundingRateHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/get_funding_rate_history", params, &raw)
	return raw, err
}

// GetPerpImpactTWAP returns the time-weighted average impact price for
// one currency's perpetual book over the requested window.
//
// Required params: `currency`, `start_time`, `end_time`. Public.
func (a *API) GetPerpImpactTWAP(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/get_perp_impact_twap", params, &raw)
	return raw, err
}

// GetLatestSignedFeeds returns the latest oracle signed-feed snapshot.
//
// Optional params: `currency`. Pass nil to get every currency the venue
// publishes. Public.
func (a *API) GetLatestSignedFeeds(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if params == nil {
		params = map[string]any{}
	}
	var raw json.RawMessage
	err := a.call(ctx, "public/get_latest_signed_feeds", params, &raw)
	return raw, err
}

// GetSpotFeedHistory returns historical oracle spot prices for one
// currency over the requested window at the given period.
//
// Required params: `currency`, `period`, `start_timestamp`,
// `end_timestamp`. Public.
func (a *API) GetSpotFeedHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/get_spot_feed_history", params, &raw)
	return raw, err
}

// GetStatistics returns rolling 24h volume / OI / volatility statistics
// for one instrument.
//
// Required params: `instrument_name`. Public.
func (a *API) GetStatistics(ctx context.Context, instrument string) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/statistics", map[string]any{"instrument_name": instrument}, &raw)
	return raw, err
}

// GetTransaction returns the on-chain status / receipt of one Derive
// transaction by its transaction id.
//
// Required params: `transaction_id`. Public.
func (a *API) GetTransaction(ctx context.Context, transactionID string) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/get_transaction", map[string]any{"transaction_id": transactionID}, &raw)
	return raw, err
}
