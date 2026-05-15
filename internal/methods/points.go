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

// GetPoints returns one wallet's points record for one program.
// Public.
func (a *API) GetPoints(ctx context.Context, q types.PointsQuery) (types.PointsRecord, error) {
	var resp types.PointsRecord
	if err := a.call(ctx, "public/get_points", map[string]any{
		"program": q.Program,
		"wallet":  q.Wallet,
	}, &resp); err != nil {
		return types.PointsRecord{}, err
	}
	return resp, nil
}

// GetPointsLeaderboard returns one page (up to 500 entries) of the
// points leaderboard for one program, ordered by points desc.
// Public.
//
// Pass `page` 1-indexed; the response carries the total number of
// pages.
func (a *API) GetPointsLeaderboard(ctx context.Context, q types.PointsLeaderboardQuery) (types.PointsLeaderboard, error) {
	params := map[string]any{"program": q.Program}
	if q.Page > 0 {
		params["page"] = q.Page
	}
	var resp types.PointsLeaderboard
	if err := a.call(ctx, "public/get_points_leaderboard", params, &resp); err != nil {
		return types.PointsLeaderboard{}, err
	}
	return resp, nil
}

// GetAllPoints returns the program-wide points snapshot for one
// program: aggregate notional volume, user count, and per-wallet
// points map. Public.
//
// The per-wallet `points` map is preserved as raw JSON because the
// inner schema varies per program; decode further at the call site.
func (a *API) GetAllPoints(ctx context.Context, q types.AllPointsQuery) (types.AllPointsResult, error) {
	var resp types.AllPointsResult
	if err := a.call(ctx, "public/get_all_points", map[string]any{"program": q.Program}, &resp); err != nil {
		return types.AllPointsResult{}, err
	}
	return resp, nil
}
