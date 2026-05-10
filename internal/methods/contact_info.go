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

// GetContactInfo lists every contact-info record on the configured
// signer's wallet, optionally filtered by `contact_type`. Private.
//
// Pass an empty `contactType` to return every record.
func (a *API) GetContactInfo(ctx context.Context, contactType string) ([]types.Contact, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"wallet": a.Signer.Owner().Hex(),
	}
	if contactType != "" {
		params["contact_type"] = contactType
	}
	var resp struct {
		Contacts []types.Contact `json:"contacts"`
	}
	if err := a.call(ctx, "private/get_contact_info", params, &resp); err != nil {
		return nil, err
	}
	return resp.Contacts, nil
}

// CreateContactInfo registers a new contact-info record on the
// configured signer's wallet. Private.
//
// The wallet param is sourced from the configured signer; pass the
// new contact's `contact_type` (e.g. "email", "telegram") and
// `contact_value`.
func (a *API) CreateContactInfo(ctx context.Context, contactType, contactValue string) (*types.Contact, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"wallet":        a.Signer.Owner().Hex(),
		"contact_type":  contactType,
		"contact_value": contactValue,
	}
	var resp struct {
		ContactInfo types.Contact `json:"contact_info"`
	}
	if err := a.call(ctx, "private/create_contact_info", params, &resp); err != nil {
		return nil, err
	}
	return &resp.ContactInfo, nil
}
