// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file collects miscellaneous single-input query DTOs that don't
// belong to a larger grouping (orders by id, referrals, vaults,
// contact info, statistics, points, MMP, maker programs).
package types

// OrderQuery parameterises private/get_order.
type OrderQuery struct {
	// OrderID is the engine-assigned id.
	OrderID string
}

// ContactInfoQuery parameterises private/get_contact_info.
type ContactInfoQuery struct {
	// ContactType filters by contact kind (e.g. "email"). Empty
	// returns every kind for the caller.
	ContactType string
}

// MMPConfigQuery parameterises private/get_mmp_config.
type MMPConfigQuery struct {
	// Currency selects the MMP bucket. Empty returns every bucket
	// the caller has configured.
	Currency string
}

// ReferralCodeQuery parameterises public/get_referral_code.
type ReferralCodeQuery struct {
	// Wallet is the Ethereum wallet address to look up.
	Wallet string
}

// InviteCodeQuery parameterises public/get_invite_code.
type InviteCodeQuery struct {
	// Wallet is the Ethereum wallet address to look up.
	Wallet string
}

// ValidateInviteCodeQuery parameterises public/validate_invite_code.
type ValidateInviteCodeQuery struct {
	// Code is the invite code to validate.
	Code string
}

// VaultBalancesQuery parameterises public/get_vault_balances.
type VaultBalancesQuery struct {
	// Wallet is the user wallet whose vault balances are queried.
	Wallet string
	// SmartContractOwner is the on-chain owner of the vault contract.
	SmartContractOwner string
}

// VaultRatesQuery parameterises public/get_vault_rates.
type VaultRatesQuery struct {
	// VaultType selects the vault category.
	VaultType string
}

// UserStatisticsQuery parameterises public/user_statistics.
type UserStatisticsQuery struct {
	// Wallet is the wallet address whose statistics are queried.
	Wallet string
}

// AllStatisticsQuery parameterises public/all_statistics.
type AllStatisticsQuery struct {
	// EndTime is the Unix-seconds anchor for the rolling window.
	// Zero defers to "now".
	EndTime int64
}

// AllUserStatisticsQuery parameterises public/all_user_statistics.
type AllUserStatisticsQuery struct {
	// EndTimeSec is the Unix-seconds anchor for the rolling window.
	// Zero defers to "now".
	EndTimeSec int64
}

// PointsQuery parameterises public/get_points.
type PointsQuery struct {
	// Program is the points program identifier.
	Program string
	// Wallet is the wallet whose points are queried.
	Wallet string
}

// PointsLeaderboardQuery parameterises public/get_points_leaderboard.
type PointsLeaderboardQuery struct {
	// Program is the points program identifier.
	Program string
	// Page is the leaderboard page number (1-indexed).
	Page int
}

// AllPointsQuery parameterises public/get_all_points.
type AllPointsQuery struct {
	// Program selects the points program.
	Program string
}

// MakerProgramScoresQuery parameterises public/get_maker_program_scores.
type MakerProgramScoresQuery struct {
	// ProgramName identifies the maker program.
	ProgramName string
	// EpochStartTimestamp is the Unix-seconds start of the epoch
	// whose scores are being read.
	EpochStartTimestamp int64
}
