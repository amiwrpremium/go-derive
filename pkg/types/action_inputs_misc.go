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

// Validate enforces Label is populated.
func (in ChangeSubaccountLabelInput) Validate() error {
	if in.Label == "" {
		return invalidParam("label", "required")
	}
	return nil
}

// CreateContactInfoInput parameterises private/create_contact_info.
type CreateContactInfoInput struct {
	// ContactType is the kind of contact info (email, telegram, etc.).
	ContactType string
	// ContactValue is the contact value to register.
	ContactValue string
}

// Validate enforces both fields are populated.
func (in CreateContactInfoInput) Validate() error {
	if in.ContactType == "" {
		return invalidParam("contact_type", "required")
	}
	if in.ContactValue == "" {
		return invalidParam("contact_value", "required")
	}
	return nil
}

// UpdateContactInfoInput parameterises private/update_contact_info.
type UpdateContactInfoInput struct {
	// ContactID is the server-side id of the contact record to update.
	ContactID int64
	// NewValue is the replacement value for the contact.
	NewValue string
}

// Validate enforces both fields are populated.
func (in UpdateContactInfoInput) Validate() error {
	if in.ContactID == 0 {
		return invalidParam("contact_id", "required")
	}
	if in.NewValue == "" {
		return invalidParam("contact_value", "required")
	}
	return nil
}

// DeleteContactInfoInput parameterises private/delete_contact_info.
type DeleteContactInfoInput struct {
	// ContactID is the server-side id of the contact record to delete.
	ContactID int64
}

// Validate enforces ContactID is populated.
func (in DeleteContactInfoInput) Validate() error {
	if in.ContactID == 0 {
		return invalidParam("contact_id", "required")
	}
	return nil
}

// ResetMMPInput parameterises private/reset_mmp.
type ResetMMPInput struct {
	// Currency selects the MMP bucket to reset.
	Currency string
}

// Validate enforces Currency is populated.
func (in ResetMMPInput) Validate() error {
	if in.Currency == "" {
		return invalidParam("currency", "required")
	}
	return nil
}

// SetCancelOnDisconnectInput parameterises
// private/set_cancel_on_disconnect.
type SetCancelOnDisconnectInput struct {
	// Enabled turns cancel-on-disconnect on or off for the connection.
	Enabled bool
}
