// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shapes for the public maker-program and
// referral-performance endpoints: `public/get_maker_programs`,
// `public/get_maker_program_scores`, and
// `public/get_referral_performance`.
package types

// MakerProgram is one entry in `public/get_maker_programs` and the
// `program` field of [MakerProgramScore]. It describes one
// maker-incentive epoch: which assets and currencies it covers, the
// minimum notional for eligibility, the reward tokens, and the
// epoch's start / end timestamps.
//
// The shape mirrors `ProgramResponseSchema` in Derive's v2.2 OpenAPI
// spec.
type MakerProgram struct {
	// Name is the program name (e.g. "options-maker-q1").
	Name string `json:"name"`
	// AssetTypes is the list of instrument-asset types covered by the
	// program (e.g. "option", "perp").
	AssetTypes []string `json:"asset_types"`
	// Currencies is the list of underlying currencies covered (e.g.
	// "ETH", "BTC").
	Currencies []string `json:"currencies"`
	// MinNotional is the minimum dollar notional a market maker must
	// quote to qualify for the program.
	MinNotional Decimal `json:"min_notional"`
	// Rewards maps reward-token symbol → total reward amount paid out
	// over the epoch.
	Rewards map[string]Decimal `json:"rewards"`
	// StartTimestamp is the epoch start (Unix seconds — the OAS does
	// not annotate the unit, but matches existing program-related
	// timestamps elsewhere).
	StartTimestamp int64 `json:"start_timestamp"`
	// EndTimestamp is the epoch end.
	EndTimestamp int64 `json:"end_timestamp"`
}

// MakerProgramScore is the response of
// `public/get_maker_program_scores`. It reports the per-wallet score
// breakdown for one maker-incentive program, plus the program totals.
//
// The shape mirrors `PublicGetMakerProgramScoresResultSchema`.
type MakerProgramScore struct {
	// Program is the program metadata (same as one entry in
	// `public/get_maker_programs`).
	Program MakerProgram `json:"program"`
	// Scores is the per-wallet score breakdown.
	Scores []MakerScoreBreakdown `json:"scores"`
	// TotalScore is the total score across all market makers for the
	// epoch.
	TotalScore Decimal `json:"total_score"`
	// TotalVolume is the total volume across all market makers for
	// the epoch.
	TotalVolume Decimal `json:"total_volume"`
}

// MakerScoreBreakdown is one wallet's per-program score, as returned
// in [MakerProgramScore.Scores].
//
// The shape mirrors `ScoreBreakdownSchema`.
type MakerScoreBreakdown struct {
	// Wallet is the market maker's wallet address.
	Wallet Address `json:"wallet"`
	// CoverageScore is the coverage component of the score.
	CoverageScore Decimal `json:"coverage_score"`
	// QualityScore is the quality component of the score.
	QualityScore Decimal `json:"quality_score"`
	// HolderBoost is a per-account multiplier for the score due to
	// holding tokens.
	HolderBoost Decimal `json:"holder_boost"`
	// Volume is the volume traded by the account for this epoch.
	Volume Decimal `json:"volume"`
	// VolumeMultiplier is the multiplier applied to the volume for
	// scoring.
	VolumeMultiplier Decimal `json:"volume_multiplier"`
	// TotalScore is the total score of the account for this program.
	TotalScore Decimal `json:"total_score"`
}

// ReferralPerformance is the response of
// `public/get_referral_performance`. It reports the headline referrer
// performance for one referral code, plus a deeply-nested breakdown
// keyed by liquidity role → currency → instrument type.
//
// The shape mirrors `PublicGetReferralPerformanceResultSchema`.
type ReferralPerformance struct {
	// ReferralCode is the referral code the metrics are scoped to.
	ReferralCode string `json:"referral_code"`
	// FeeSharePercentage is the percentage of referred fees rebated
	// to the referrer.
	FeeSharePercentage Decimal `json:"fee_share_percentage"`
	// StdrvBalance is the staked-DRV balance the referrer holds (used
	// to determine the fee-share percentage).
	StdrvBalance Decimal `json:"stdrv_balance"`
	// TotalNotionalVolume is the total referred notional volume.
	TotalNotionalVolume Decimal `json:"total_notional_volume"`
	// TotalReferredFees is the total fees paid by referred traders.
	TotalReferredFees Decimal `json:"total_referred_fees"`
	// TotalFeeRewards is the total fee rewards paid to referrers.
	TotalFeeRewards Decimal `json:"total_fee_rewards"`
	// TotalBuilderFeeCollected is the total builder fee collected via
	// the `extra_fee` field on referred orders.
	TotalBuilderFeeCollected Decimal `json:"total_builder_fee_collected"`
	// Rewards is the deeply-nested per-role / per-currency /
	// per-instrument-type breakdown. Keys at each level mirror the
	// OAS verbatim — outer = liquidity role ("maker" / "taker"),
	// middle = currency, inner = instrument type.
	Rewards map[string]map[string]map[string]ReferralBreakdown `json:"rewards"`
}

// ReferralBreakdown is one leaf of [ReferralPerformance.Rewards].
//
// The shape mirrors `ReferralPerformanceByInstrumentTypeSchema`.
type ReferralBreakdown struct {
	// NotionalVolume is the notional volume traded under this slice.
	NotionalVolume Decimal `json:"notional_volume"`
	// ReferredFee is the fees paid by the referred trader.
	ReferredFee Decimal `json:"referred_fee"`
	// FeeReward is the fee rebate paid to the referrer.
	FeeReward Decimal `json:"fee_reward"`
	// BuilderFee is the builder fee collected via `extra_fee`.
	BuilderFee Decimal `json:"builder_fee"`
	// UniqueTradersReferred is the number of unique wallets referred
	// under this slice.
	UniqueTradersReferred int64 `json:"unique_traders_referred"`
}
