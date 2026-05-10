package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestCreateContactInfo_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/create_contact_info", map[string]any{
		"contact_info": map[string]any{
			"id":             int64(7),
			"contact_type":   "email",
			"contact_value":  "alice@example.com",
			"created_at_sec": int64(1700000000),
			"updated_at_sec": int64(1700000000),
		},
	})
	got, err := api.CreateContactInfo(context.Background(), "email", "alice@example.com")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(7), got.ID)
	assert.Equal(t, "alice@example.com", got.ContactValue)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.NotEmpty(t, params["wallet"])
	assert.Equal(t, "email", params["contact_type"])
}

func TestCreateContactInfo_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.CreateContactInfo(context.Background(), "email", "x@y.z")
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
}

func TestGetContactInfo_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/get_contact_info", map[string]any{
		"contacts": []any{
			map[string]any{
				"id": int64(1), "contact_type": "email", "contact_value": "alice@example.com",
				"created_at_sec": int64(1700000000), "updated_at_sec": int64(1700000000),
			},
			map[string]any{
				"id": int64(2), "contact_type": "telegram", "contact_value": "@alice",
				"created_at_sec": int64(1700000010), "updated_at_sec": int64(1700000010),
			},
		},
	})
	got, err := api.GetContactInfo(context.Background(), "")
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "telegram", got[1].ContactType)
}

func TestGetContactInfo_AppliesTypeFilter(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/get_contact_info", map[string]any{"contacts": []any{}})
	_, err := api.GetContactInfo(context.Background(), "email")
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "email", params["contact_type"])
}

func TestUpdateContactInfo_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/update_contact_info", map[string]any{
		"contact_info": map[string]any{
			"id": int64(7), "contact_type": "email", "contact_value": "alice@new.com",
			"created_at_sec": int64(1700000000), "updated_at_sec": int64(1700000060),
		},
	})
	got, err := api.UpdateContactInfo(context.Background(), 7, "alice@new.com")
	require.NoError(t, err)
	assert.Equal(t, "alice@new.com", got.ContactValue)
	assert.Equal(t, int64(1700000060), got.UpdatedAtSec)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, float64(7), params["contact_id"])
	assert.Equal(t, "alice@new.com", params["contact_value"])
}

func TestDeleteContactInfo_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/delete_contact_info", map[string]any{
		"contact_id": int64(7),
		"deleted":    true,
	})
	id, deleted, err := api.DeleteContactInfo(context.Background(), 7)
	require.NoError(t, err)
	assert.Equal(t, int64(7), id)
	assert.True(t, deleted)
}

func TestContactInfo_AllRequireSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.GetContactInfo(context.Background(), "")
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
	_, err = api.UpdateContactInfo(context.Background(), 1, "x")
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
	_, _, err = api.DeleteContactInfo(context.Background(), 1)
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
}
