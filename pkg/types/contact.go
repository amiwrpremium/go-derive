// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the contact-info entity returned by the
// private/{create,get,update,delete}_contact_info endpoints.
package types

// Contact is one stored contact-info record on a wallet — typically
// an email or messaging handle the engine can reach the wallet's
// owner on (e.g. for liquidation alerts).
//
// The shape mirrors `ContactInfoSchema` per
// docs.derive.xyz/reference/private-{create,get}_contact_info.
type Contact struct {
	// ID is the server-side contact id (used as the handle for
	// update / delete).
	ID int64 `json:"id"`
	// ContactType identifies the channel — e.g. "email", "telegram".
	ContactType string `json:"contact_type"`
	// ContactValue is the channel-specific address (the email
	// itself, the telegram handle, etc.).
	ContactValue string `json:"contact_value"`
	// CreatedAtSec is the Unix-seconds creation timestamp.
	CreatedAtSec int64 `json:"created_at_sec"`
	// UpdatedAtSec is the Unix-seconds last-update timestamp.
	UpdatedAtSec int64 `json:"updated_at_sec"`
}
