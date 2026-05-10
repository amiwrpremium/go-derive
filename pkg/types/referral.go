// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the referral / invite-code shapes returned by
// `public/get_all_referral_codes`, `public/get_invite_code`, etc.
package types

// ReferralCodeRecord is one referral-code record. Returned by
// `public/get_all_referral_codes` (slice).
type ReferralCodeRecord struct {
	// ReferralCode is the code string itself.
	ReferralCode string `json:"referral_code"`
	// Wallet is the referrer's wallet address.
	Wallet string `json:"wallet"`
	// ReceivingWallet is the wallet that receives the rebates.
	ReceivingWallet string `json:"receiving_wallet"`
}

// InviteCode is the response of `public/get_invite_code` — the
// invite code allocated to one wallet plus its remaining-uses
// counter (`-1` for unlimited).
type InviteCode struct {
	// Code is the invite code string.
	Code string `json:"code"`
	// RemainingUses is how many more uses the code has. -1 means
	// unlimited.
	RemainingUses int64 `json:"remaining_uses"`
}
