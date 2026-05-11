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

// GetNotifications returns the wallet's notification feed. Private.
//
// Each [types.Notification.EventDetails] field is open-shaped per
// the OAS — decode against the concrete event type at the call
// site if you need it.
func (a *API) GetNotifications(ctx context.Context, q types.NotificationsQuery, page types.PageRequest) ([]types.Notification, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	params := map[string]any{}
	if a.Subaccount != 0 {
		params["subaccount_id"] = a.Subaccount
	}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	}
	if q.Status != "" {
		params["status"] = q.Status
	}
	if len(q.Types) > 0 {
		params["type"] = q.Types
	}
	addPaging(params, page)
	var resp struct {
		Notifications []types.Notification `json:"notifications"`
		Pagination    types.Page           `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_notifications", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Notifications, resp.Pagination, nil
}

// UpdateNotifications marks one or more notifications as seen or
// hidden. Private. Returns the number of notifications updated.
func (a *API) UpdateNotifications(ctx context.Context, in types.UpdateNotificationsInput) (types.UpdateNotificationsResult, error) {
	if err := a.requireSigner(); err != nil {
		return types.UpdateNotificationsResult{}, err
	}
	sub := in.SubaccountID
	if sub == 0 {
		sub = a.Subaccount
	}
	params := map[string]any{
		"subaccount_id":    sub,
		"notification_ids": in.NotificationIDs,
	}
	if in.Status != "" {
		params["status"] = in.Status
	}
	var resp types.UpdateNotificationsResult
	if err := a.call(ctx, "private/update_notifications", params, &resp); err != nil {
		return types.UpdateNotificationsResult{}, err
	}
	return resp, nil
}
