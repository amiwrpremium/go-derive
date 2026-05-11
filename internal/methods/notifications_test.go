package methods_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestGetNotifications_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/get_notifications", map[string]any{
		"notifications": []any{
			map[string]any{
				"id":             int64(42),
				"subaccount_id":  int64(7),
				"event":          "deposit_completed",
				"event_details":  map[string]any{"asset": "USDC"},
				"status":         "unseen",
				"timestamp":      int64(1700000000000),
				"transaction_id": int64(123),
				"tx_hash":        "0xabc",
			},
		},
		"pagination": map[string]any{"count": 1, "num_pages": 1},
	})
	notifs, page, err := api.GetNotifications(context.Background(), types.NotificationsQuery{}, types.PageRequest{})
	require.NoError(t, err)
	require.Len(t, notifs, 1)
	assert.Equal(t, "deposit_completed", notifs[0].Event)
	assert.JSONEq(t, `{"asset":"USDC"}`, string(notifs[0].EventDetails))
	require.NotNil(t, notifs[0].TransactionID)
	assert.Equal(t, int64(123), *notifs[0].TransactionID)
	assert.Equal(t, 1, page.Count)
}

func TestGetNotifications_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, _, err := api.GetNotifications(context.Background(), types.NotificationsQuery{}, types.PageRequest{})
	assert.True(t, errors.Is(err, derrors.ErrUnauthorized))
}

func TestGetNotifications_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleError("private/get_notifications", boom)
	_, _, err := api.GetNotifications(context.Background(), types.NotificationsQuery{}, types.PageRequest{})
	assert.ErrorAs(t, err, new(*derrors.APIError))
}

func TestUpdateNotifications_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/update_notifications", map[string]any{"updated_count": int64(3)})
	got, err := api.UpdateNotifications(context.Background(), types.UpdateNotificationsInput{
		NotificationIDs: []int64{1, 2, 3},
		Status:          "seen",
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(3), got.UpdatedCount)
}

func TestUpdateNotifications_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.UpdateNotifications(context.Background(), types.UpdateNotificationsInput{})
	assert.True(t, errors.Is(err, derrors.ErrUnauthorized))
}

func TestUpdateNotifications_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleError("private/update_notifications", boom)
	_, err := api.UpdateNotifications(context.Background(), types.UpdateNotificationsInput{})
	assert.ErrorAs(t, err, new(*derrors.APIError))
}
