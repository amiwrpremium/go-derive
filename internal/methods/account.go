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

// GetAccount returns wallet-level account information for the
// configured signer's owner. Private.
//
// Reports the wallet's subaccount ids, the kill-switch state, the
// per-WS-budget TPS limits, the fee schedule, and the wallet's
// referral code.
func (a *API) GetAccount(ctx context.Context) (*types.AccountResult, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	params := map[string]any{"wallet": a.Signer.Owner().Hex()}
	var resp types.AccountResult
	if err := a.call(ctx, "private/get_account", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAllPortfolios returns the per-subaccount portfolio snapshot for
// every subaccount owned by the configured signer's wallet. Private.
//
// Each entry is a full snapshot — collateral, positions, open orders,
// and the engine-side margin breakdown — for one subaccount.
//
// To query a different wallet, populate `wallet` directly via the
// raw transport.
func (a *API) GetAllPortfolios(ctx context.Context) ([]types.Portfolio, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	params := map[string]any{"wallet": a.Signer.Owner().Hex()}
	var resp []types.Portfolio
	if err := a.call(ctx, "private/get_all_portfolios", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetMargin returns the live margin breakdown for the configured
// subaccount. Private.
//
// Pre/post-margin values report the impact of the simulated trade
// implied by `private/get_margin`'s `simulated_*` parameters; with
// no simulation the pre and post columns are equal.
func (a *API) GetMargin(ctx context.Context) (*types.MarginResult, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var resp types.MarginResult
	if err := a.call(ctx, "private/get_margin", map[string]any{"subaccount_id": a.Subaccount}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPublicMargin runs Derive's risk-engine margin calculation
// against the user-supplied collaterals and positions, returning
// the resulting margin requirement. Public — no signer required.
//
// Pass the same shape `private/get_margin` accepts as `params`.
// Required keys per the OAS: `margin_type` ("PM" / "PM2" / "SM"),
// `market`, `simulated_collaterals`, `simulated_positions`. Optional:
// `simulated_collateral_changes`, `simulated_position_changes`.
func (a *API) GetPublicMargin(ctx context.Context, params map[string]any) (*types.MarginResult, error) {
	var resp types.MarginResult
	if err := a.call(ctx, "public/get_margin", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
