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

// GetMakerPrograms returns the list of currently-active maker
// incentive programs. Public.
//
// Each program has its own epoch (start / end timestamps), the
// instrument-asset types and currencies it covers, the minimum
// dollar notional for eligibility, and the reward token amounts paid
// out over the epoch.
func (a *API) GetMakerPrograms(ctx context.Context) ([]types.MakerProgram, error) {
	var resp []types.MakerProgram
	if err := a.call(ctx, "public/get_maker_programs", map[string]any{}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetMakerProgramScores returns the per-wallet score breakdown for
// one maker incentive program at one epoch. Public.
//
// Required `params`: `program_name`, `epoch_start_timestamp`. The
// response carries the program metadata alongside the per-wallet
// breakdown and the program-wide totals.
func (a *API) GetMakerProgramScores(ctx context.Context, params map[string]any) (*types.MakerProgramScore, error) {
	var resp types.MakerProgramScore
	if err := a.call(ctx, "public/get_maker_program_scores", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDetailedMakerSnapshotHistory returns the per-quote maker
// snapshots for one program / epoch, paginated. Public.
//
// Required `params`: `program_name`, `epoch_start_timestamp`.
// Optional: `wallet` (filter to one maker), `page`, `page_size`.
func (a *API) GetDetailedMakerSnapshotHistory(ctx context.Context, params map[string]any) (*types.DetailedMakerSnapshotHistory, error) {
	if params == nil {
		params = map[string]any{}
	}
	var resp types.DetailedMakerSnapshotHistory
	if err := a.call(ctx, "public/get_detailed_maker_snapshot_history", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetReferralPerformance returns the headline referrer performance
// for one referral code over the requested time window, plus a
// deeply-nested per-role / per-currency / per-instrument-type
// breakdown. Public.
//
// Required `params`: `start_ms`, `end_ms`. Optional: `referral_code`,
// `wallet`.
func (a *API) GetReferralPerformance(ctx context.Context, params map[string]any) (*types.ReferralPerformance, error) {
	var resp types.ReferralPerformance
	if err := a.call(ctx, "public/get_referral_performance", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
