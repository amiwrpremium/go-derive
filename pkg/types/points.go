// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the points / leaderboard shapes returned by
// `public/get_points`, `public/get_all_points`, and
// `public/get_points_leaderboard`.
package types

import "encoding/json"

// PointsRecord is one user's points record for a given program.
// Returned by `public/get_points`.
//
// `percent_share_of_points`, `total_users`, and `user_rank` are
// documented as deprecated on the wire — the SDK still surfaces them
// for completeness; new callers should rely on the leaderboard
// endpoint instead.
type PointsRecord struct {
	// Flag is a server-side flag for special treatment of the user.
	Flag string `json:"flag,omitempty"`
	// Parent is the referrer wallet address.
	Parent string `json:"parent,omitempty"`
	// PercentShareOfPoints is deprecated — kept for completeness.
	PercentShareOfPoints Decimal `json:"percent_share_of_points,omitempty"`
	// TotalNotionalVolume is the user's $ notional volume traded
	// in the program.
	TotalNotionalVolume Decimal `json:"total_notional_volume,omitempty"`
	// TotalUsers is deprecated.
	TotalUsers int64 `json:"total_users,omitempty"`
	// UserRank is deprecated.
	UserRank int64 `json:"user_rank,omitempty"`
	// Points is the per-category points map (e.g. category →
	// points value). Kept as raw JSON because the inner shape
	// varies per program.
	Points json.RawMessage `json:"points"`
}

// AllPointsResult is the response of `public/get_all_points`. Each
// entry in the inner `points` map is keyed by wallet — kept as raw
// JSON because the inner map's value shape varies per program.
type AllPointsResult struct {
	// TotalNotionalVolume is the program-wide $ notional volume.
	TotalNotionalVolume Decimal `json:"total_notional_volume"`
	// TotalUsers is the count of distinct users in the program.
	TotalUsers int64 `json:"total_users"`
	// Points is the wallet → per-user points map.
	Points json.RawMessage `json:"points"`
}

// LeaderboardEntry is one row of the `public/get_points_leaderboard`
// response.
type LeaderboardEntry struct {
	// Rank is the user's leaderboard rank.
	Rank int64 `json:"rank"`
	// Wallet is the user's wallet address.
	Wallet string `json:"wallet"`
	// Points is the user's total points in the program.
	Points Decimal `json:"points"`
	// PercentShareOfPoints is the user's share of the program's
	// total points (e.g. "0.025" for 2.5 %).
	PercentShareOfPoints Decimal `json:"percent_share_of_points"`
	// TotalVolume is the user's $ notional volume.
	TotalVolume Decimal `json:"total_volume"`
	// Flag is a server-side flag for special treatment of the user.
	Flag string `json:"flag,omitempty"`
	// Parent is the referrer wallet address.
	Parent string `json:"parent,omitempty"`
}

// PointsLeaderboard is the response of `public/get_points_leaderboard`.
type PointsLeaderboard struct {
	// Pages is the total number of pages in the leaderboard.
	Pages int64 `json:"pages"`
	// TotalUsers is the program-wide user count.
	TotalUsers int64 `json:"total_users"`
	// Leaderboard is the slice of leaderboard rows for the requested
	// page (up to 500 entries).
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
}
