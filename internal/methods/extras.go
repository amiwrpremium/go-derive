// Package methods.
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

// GetPublicMargin runs Derive's risk-engine margin calculation against a
// user-supplied set of `simulated_collaterals` and `simulated_positions`,
// returning the resulting margin requirement.
//
// Required params: `simulated_collaterals`, `simulated_positions`,
// `margin_type` ("PM" / "PM2" / "SM"). Public — no signer required.
func (a *API) GetPublicMargin(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/get_margin", params, &raw)
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

// GetPublicOptionSettlementHistory returns the network-wide option
// settlement history.
//
// Optional params: pagination. Pass nil for the default range. Public.
func (a *API) GetPublicOptionSettlementHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if params == nil {
		params = map[string]any{}
	}
	var raw json.RawMessage
	err := a.call(ctx, "public/get_option_settlement_history", params, &raw)
	return raw, err
}

// ---------------------------------------------------------------------------
// Private read methods
// ---------------------------------------------------------------------------

// GetAccount returns wallet-level account information for the signer.
//
// No params. Private.
func (a *API) GetAccount(ctx context.Context) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	params := map[string]any{"wallet": a.Signer.Owner().Hex()}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_account", params, &raw)
	return raw, err
}

// GetMargin returns the live margin breakdown for the configured
// subaccount.
//
// No params. Private.
func (a *API) GetMargin(ctx context.Context) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_margin", map[string]any{"subaccount_id": a.Subaccount}, &raw)
	return raw, err
}

// GetFundingHistory returns funding payments received / paid by the
// configured subaccount.
//
// Optional params: `start_timestamp`, `end_timestamp`, `instrument_name`,
// pagination. Private.
func (a *API) GetFundingHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	err := a.call(ctx, "private/get_funding_history", params, &raw)
	return raw, err
}

// GetLiquidationHistory returns the configured subaccount's past
// liquidation events.
//
// Optional params: `start_timestamp`, `end_timestamp`, pagination. Private.
func (a *API) GetLiquidationHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	err := a.call(ctx, "private/get_liquidation_history", params, &raw)
	return raw, err
}

// GetOptionSettlementHistory returns the configured subaccount's past
// option-settlement events.
//
// Optional params: `start_timestamp`, `end_timestamp`, pagination. Private.
func (a *API) GetOptionSettlementHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	err := a.call(ctx, "private/get_option_settlement_history", params, &raw)
	return raw, err
}

// GetSubaccountValueHistory returns the equity-curve series for the
// configured subaccount.
//
// Required params: `period`, `start_timestamp`, `end_timestamp`. Private.
func (a *API) GetSubaccountValueHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	err := a.call(ctx, "private/get_subaccount_value_history", params, &raw)
	return raw, err
}

// GetERC20TransferHistory returns deposit / withdrawal-style ERC-20
// transfers attributed to the configured subaccount.
func (a *API) GetERC20TransferHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	err := a.call(ctx, "private/get_erc20_transfer_history", params, &raw)
	return raw, err
}

// GetInterestHistory returns the configured subaccount's interest charges
// and rebates.
func (a *API) GetInterestHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	err := a.call(ctx, "private/get_interest_history", params, &raw)
	return raw, err
}

// ExpiredAndCancelledHistory returns the configured subaccount's expired
// and cancelled orders.
func (a *API) ExpiredAndCancelledHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
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
	err := a.call(ctx, "private/expired_and_cancelled_history", params, &raw)
	return raw, err
}

// GetMMPConfig returns the active market-maker-protection config for the
// configured subaccount.
func (a *API) GetMMPConfig(ctx context.Context) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_mmp_config", map[string]any{"subaccount_id": a.Subaccount}, &raw)
	return raw, err
}

// GetNotifications returns the wallet's notification feed.
func (a *API) GetNotifications(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok && a.Subaccount != 0 {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_notifications", params, &raw)
	return raw, err
}

// UpdateNotifications marks one or more notifications as seen / dismissed.
//
// Required params: `notification_ids` ([]int) and `status`. Private.
func (a *API) UpdateNotifications(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/update_notifications", params, &raw)
	return raw, err
}

// ---------------------------------------------------------------------------
// Private write methods
// ---------------------------------------------------------------------------

// Replace cancels one outstanding order and submits a replacement in a
// single round trip — the standard maker pattern for re-pricing without a
// race against the matching engine.
//
// Params should include `order_id_to_cancel` and the same fields PlaceOrder
// would take. The full param shape is documented at docs.derive.xyz.
//
// Private; requires signer + subaccount.
func (a *API) Replace(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/replace", params, &raw)
	return raw, err
}

// OrderDebug previews an order without submitting it — Derive returns the
// validated request and any synthetic fees / margin impacts the engine
// computes. Use this to sanity-check signed payloads in CI.
//
// Params mirror PlaceOrder. Private.
func (a *API) OrderDebug(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/order_debug", params, &raw)
	return raw, err
}

// CancelByNonce cancels an order by the nonce on its signed payload —
// useful when the caller has not received the order id back yet.
//
// Required params: `instrument_name`, `nonce`, `wallet`. Private.
func (a *API) CancelByNonce(ctx context.Context, instrument string, nonce uint64) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"instrument_name": instrument,
		"nonce":           nonce,
		"wallet":          a.Signer.Owner().Hex(),
		"subaccount_id":   a.Subaccount,
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/cancel_by_nonce", params, &raw)
	return raw, err
}

// SetCancelOnDisconnect arms or disarms the kill-switch that cancels every
// open order on the wallet if the WebSocket session disconnects.
//
// Pass enabled=true to arm; false to disarm. Private.
func (a *API) SetCancelOnDisconnect(ctx context.Context, enabled bool) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"wallet":  a.Signer.Owner().Hex(),
		"enabled": enabled,
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/set_cancel_on_disconnect", params, &raw)
	return raw, err
}

// ChangeSubaccountLabel sets the human-readable label on the configured
// subaccount.
func (a *API) ChangeSubaccountLabel(ctx context.Context, label string) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"label":         label,
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/change_subaccount_label", params, &raw)
	return raw, err
}
