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
// Optional `params`: `status` ("unseen" / "seen" / "hidden"),
// `type` ([]string), `page`, `page_size`. The configured
// subaccount is threaded through automatically when set and not
// already present in `params`.
//
// Each [types.Notification.EventDetails] field is open-shaped per
// the OAS — decode against the concrete event type at the call
// site if you need it.
func (a *API) GetNotifications(ctx context.Context, params map[string]any) ([]types.Notification, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok && a.Subaccount != 0 {
		params["subaccount_id"] = a.Subaccount
	}
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
// hidden. Private.
//
// Required `params`: `notification_ids` ([]int) and `status`
// ("seen" or "hidden"). Optional: `subaccount_id`. Returns the
// number of notifications updated.
func (a *API) UpdateNotifications(ctx context.Context, params map[string]any) (*types.UpdateNotificationsResult, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	var resp types.UpdateNotificationsResult
	if err := a.call(ctx, "private/update_notifications", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
