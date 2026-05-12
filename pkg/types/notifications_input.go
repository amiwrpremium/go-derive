// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds input DTOs for the notification, RFQ-listing,
// vault-share and program-history endpoints.
package types

// NotificationsQuery narrows a paginated
// `private/get_notifications` request. Wallet, when non-empty,
// takes precedence over the configured subaccount per the engine
// contract.
type NotificationsQuery struct {
	// Status filters by notification status: "unseen", "seen",
	// "hidden".
	Status string
	// Types narrows by notification type code(s).
	Types []string
	// Wallet overrides the configured subaccount when set.
	Wallet string
}

// UpdateNotificationsInput parameterises
// `private/update_notifications`. NotificationIDs must be
// non-empty and Status must be "seen" or "hidden" (the engine
// rejects "unseen").
type UpdateNotificationsInput struct {
	// SubaccountID is the subaccount the notifications belong to.
	// Zero defaults to the client-configured subaccount.
	SubaccountID int64
	// NotificationIDs is the list of notification ids to update.
	NotificationIDs []int64
	// Status is the target status: "seen" or "hidden".
	Status string
}

// Validate performs schema-level checks on the receiver.
func (in UpdateNotificationsInput) Validate() error {
	if len(in.NotificationIDs) == 0 {
		return invalidParam("notification_ids", "must be non-empty")
	}
	switch in.Status {
	case "seen", "hidden":
		return nil
	}
	return invalidParam("status", "must be one of seen, hidden")
}

// RFQsQuery narrows a paginated `private/get_rfqs` request to one
// RFQ, one status, or a `[from, to]` last-update window in
// milliseconds since the Unix epoch.
type RFQsQuery struct {
	HistoryWindow
	// RFQID restricts the result to one RFQ id.
	RFQID string
	// Status restricts the result by RFQ status: "open", "filled",
	// "cancelled", "expired".
	Status string
}

// QuotesQuery narrows a paginated `private/get_quotes` request.
// Same shape as [RFQsQuery] plus an optional QuoteID filter.
type QuotesQuery struct {
	HistoryWindow
	// RFQID restricts the result to one RFQ id.
	RFQID string
	// QuoteID restricts the result to one quote id.
	QuoteID string
	// Status restricts the result by quote status: "open",
	// "filled", "cancelled", "expired".
	Status string
}

// PollQuotesQuery is the long-poll variant of [QuotesQuery] for
// `private/poll_quotes`. Identical fields.
type PollQuotesQuery = QuotesQuery

// VaultShareQuery narrows a paginated `public/get_vault_share`
// request. Timestamps are in seconds since the Unix epoch — note
// the unit difference from the milliseconds-based history methods.
type VaultShareQuery struct {
	// VaultName is the vault token's name. Required.
	VaultName string
	// FromSec is the inclusive start of the window in seconds.
	// Required.
	FromSec int64
	// ToSec is the inclusive end of the window in seconds.
	// Required.
	ToSec int64
}

// Validate performs schema-level checks on the receiver.
func (q VaultShareQuery) Validate() error {
	if q.VaultName == "" {
		return invalidParam("vault_name", "required")
	}
	if q.FromSec <= 0 {
		return invalidParam("from_timestamp_sec", "required")
	}
	if q.ToSec <= 0 {
		return invalidParam("to_timestamp_sec", "required")
	}
	return nil
}

// DetailedMakerSnapshotHistoryQuery parameterises
// `public/get_detailed_maker_snapshot_history`, returning per-quote
// maker snapshots for one program / epoch.
type DetailedMakerSnapshotHistoryQuery struct {
	// ProgramName is the maker program. Required.
	ProgramName string
	// EpochStartTimestamp identifies the epoch (Unix milliseconds).
	// Required.
	EpochStartTimestamp int64
	// Wallet restricts the result to one maker.
	Wallet string
}

// Validate performs schema-level checks on the receiver.
func (q DetailedMakerSnapshotHistoryQuery) Validate() error {
	if q.ProgramName == "" {
		return invalidParam("program_name", "required")
	}
	if q.EpochStartTimestamp == 0 {
		return invalidParam("epoch_start_timestamp", "required")
	}
	if q.Wallet == "" {
		return invalidParam("wallet", "required")
	}
	return nil
}

// ReferralPerformanceQuery parameterises
// `public/get_referral_performance`. StartMs / EndMs are required;
// either ReferralCode or Wallet (or both) can be set to scope the
// lookup.
type ReferralPerformanceQuery struct {
	// StartMs is the start of the window in milliseconds.
	// Required.
	StartMs int64
	// EndMs is the end of the window in milliseconds. Required.
	EndMs int64
	// ReferralCode restricts the result to one referral code.
	ReferralCode string
	// Wallet restricts the result to one referrer wallet.
	Wallet string
}

// Validate performs schema-level checks on the receiver.
func (q ReferralPerformanceQuery) Validate() error {
	if q.StartMs <= 0 {
		return invalidParam("start_ms", "required")
	}
	if q.EndMs <= 0 {
		return invalidParam("end_ms", "required")
	}
	return nil
}
