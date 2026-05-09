// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the notification-feed shapes returned by
// `private/get_notifications` and `private/update_notifications`.
package types

import "encoding/json"

// Notification is one entry on the wallet's notification feed.
//
// `event_details` carries event-specific context whose shape varies
// by `event` discriminator; the OAS declares the field as an open
// `additionalProperties:{}` object, so the SDK keeps it as a raw
// payload — decode at the call site against the latest event-type
// docs.
//
// The shape mirrors `NotificationResponseSchema` in Derive's v2.2
// OpenAPI spec.
type Notification struct {
	// ID is the unique notification id.
	ID int64 `json:"id"`
	// SubaccountID is the subaccount the notification is attributed
	// to.
	SubaccountID int64 `json:"subaccount_id"`
	// Event is the discriminator string that names the kind of event
	// (e.g. "deposit_completed", "liquidation").
	Event string `json:"event"`
	// EventDetails is the per-event payload. Open shape per the OAS.
	EventDetails json.RawMessage `json:"event_details"`
	// Status is the read state ("unseen", "seen", "hidden").
	Status string `json:"status"`
	// Timestamp is when the notification was emitted (millisecond
	// Unix epoch).
	Timestamp MillisTime `json:"timestamp"`
	// TransactionID is the related Derive transaction id, if any.
	// The wire field is nullable; an absent value decodes to nil.
	TransactionID *int64 `json:"transaction_id,omitempty"`
	// TxHash is the related on-chain transaction hash, if any. The
	// wire field is nullable; an absent value decodes to an empty
	// string.
	TxHash string `json:"tx_hash,omitempty"`
}

// UpdateNotificationsResult is the response of
// `private/update_notifications`. It reports how many notifications
// were marked as seen / hidden by the call.
type UpdateNotificationsResult struct {
	// UpdatedCount is the number of notifications updated.
	UpdatedCount int64 `json:"updated_count"`
}
