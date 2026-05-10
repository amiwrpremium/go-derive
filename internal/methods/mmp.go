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

// SetMMPConfig configures market-maker protection for a currency. Private.
func (a *API) SetMMPConfig(ctx context.Context, cfg types.MMPConfig) error {
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	params := map[string]any{
		"subaccount_id":   a.Subaccount,
		"currency":        cfg.Currency,
		"mmp_frozen_time": cfg.MMPFrozenTimeMs,
		"mmp_interval":    cfg.MMPIntervalMs,
	}
	if cfg.MMPAmountLimit != "" {
		params["mmp_amount_limit"] = cfg.MMPAmountLimit
	}
	if cfg.MMPDeltaLimit != "" {
		params["mmp_delta_limit"] = cfg.MMPDeltaLimit
	}
	return a.call(ctx, "private/set_mmp_config", params, nil)
}

// ResetMMP unfreezes the subaccount's MMP for a currency. Private.
func (a *API) ResetMMP(ctx context.Context, currency string) error {
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	return a.call(ctx, "private/reset_mmp", map[string]any{
		"subaccount_id": a.Subaccount,
		"currency":      currency,
	}, nil)
}

// GetMMPConfig returns the active market-maker-protection rules
// for the configured subaccount, optionally filtered to a single
// currency. Private.
//
// Pass an empty currency to return every rule. The response is one
// record per (subaccount, currency) pair.
func (a *API) GetMMPConfig(ctx context.Context, currency string) ([]types.MMPConfigResult, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if currency != "" {
		params["currency"] = currency
	}
	var resp []types.MMPConfigResult
	if err := a.call(ctx, "private/get_mmp_config", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
