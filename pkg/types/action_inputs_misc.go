// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file collects small action-input DTOs that don't belong to a
// larger grouping (subaccount label, contact info CRUD, MMP reset,
// cancel-on-disconnect).
package types

// ChangeSubaccountLabelInput parameterises private/change_subaccount_label.
type ChangeSubaccountLabelInput struct {
	// Label is the new label to apply.
	Label string
}

// CreateContactInfoInput parameterises private/create_contact_info.
type CreateContactInfoInput struct {
	// ContactType is the kind of contact info (email, telegram, etc.).
	ContactType string
	// ContactValue is the contact value to register.
	ContactValue string
}

// UpdateContactInfoInput parameterises private/update_contact_info.
type UpdateContactInfoInput struct {
	// ContactID is the server-side id of the contact record to update.
	ContactID int64
	// NewValue is the replacement value for the contact.
	NewValue string
}

// DeleteContactInfoInput parameterises private/delete_contact_info.
type DeleteContactInfoInput struct {
	// ContactID is the server-side id of the contact record to delete.
	ContactID int64
}

// ResetMMPInput parameterises private/reset_mmp.
type ResetMMPInput struct {
	// Currency selects the MMP bucket to reset.
	Currency string
}

// SetCancelOnDisconnectInput parameterises
// private/set_cancel_on_disconnect.
type SetCancelOnDisconnectInput struct {
	// Enabled turns cancel-on-disconnect on or off for the connection.
	Enabled bool
}
