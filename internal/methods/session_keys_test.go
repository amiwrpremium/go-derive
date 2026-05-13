package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestSessionKeys_Success(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/session_keys", map[string]any{
		"public_session_keys": []any{
			map[string]any{
				"public_session_key": "0x1111111111111111111111111111111111111111",
				"label":              "primary",
				"scope":              "account",
				"expiry_sec":         int64(1900000000),
				"registered_sec":     int64(1700000000),
				"ip_whitelist":       []any{},
			},
			map[string]any{
				"public_session_key": "0x2222222222222222222222222222222222222222",
				"label":              "ro",
				"scope":              "read_only",
				"expiry_sec":         int64(1900000000),
				"registered_sec":     int64(1700000001),
				"ip_whitelist":       []any{"203.0.113.10"},
			},
		},
	})
	keys, err := api.SessionKeys(context.Background())
	require.NoError(t, err)
	require.Len(t, keys, 2)
	assert.Equal(t, "primary", keys[0].Label)
	assert.Equal(t, "account", keys[0].Scope)
	assert.Equal(t, int64(1900000000), keys[0].ExpirySec)
	assert.Equal(t, []string{"203.0.113.10"}, keys[1].IPWhitelist)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.NotEmpty(t, params["wallet"])
}

func TestSessionKeys_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.SessionKeys(context.Background())
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
}

func TestEditSessionKey_OnlySetsProvidedFields(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/edit_session_key", map[string]any{
		"public_session_key": "0x1111111111111111111111111111111111111111",
		"label":              "renamed",
		"scope":              "account",
		"expiry_sec":         int64(1900000000),
		"registered_sec":     int64(1700000000),
		"ip_whitelist":       []any{},
	})
	newLabel := "renamed"
	in := types.EditSessionKeyInput{
		PublicSessionKey: types.MustAddress("0x1111111111111111111111111111111111111111"),
		Label:            &newLabel,
	}
	got, err := api.EditSessionKey(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, "renamed", got.Label)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "renamed", params["label"])
	_, hasDisable := params["disable"]
	assert.False(t, hasDisable, "disable should be omitted when false")
	_, hasIP := params["ip_whitelist"]
	assert.False(t, hasIP, "ip_whitelist should be omitted when nil")
}

func TestEditSessionKey_SetsDisableAndClearableIPList(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/edit_session_key", map[string]any{
		"public_session_key": "0x1111111111111111111111111111111111111111",
		"label":              "x",
		"scope":              "account",
		"expiry_sec":         int64(1900000000),
		"registered_sec":     int64(1700000000),
		"ip_whitelist":       []any{},
	})
	clear := []string{}
	in := types.EditSessionKeyInput{
		PublicSessionKey: types.MustAddress("0x1111111111111111111111111111111111111111"),
		Disable:          true,
		IPWhitelist:      &clear,
	}
	_, err := api.EditSessionKey(context.Background(), in)
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, true, params["disable"])
	// An explicit empty slice marshals to an empty JSON array — server
	// reads that as "clear the allow-list".
	assert.NotNil(t, params["ip_whitelist"])
}

func TestRegisterScopedSessionKey_Success(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/register_scoped_session_key", map[string]any{
		"public_session_key": "0x3333333333333333333333333333333333333333",
		"label":              "scoped",
		"scope":              "read_only",
		"expiry_sec":         int64(1900000000),
		"ip_whitelist":       []any{},
		"transaction_id":     "txn-42",
	})
	in := types.RegisterScopedSessionKeyInput{
		PublicSessionKey: types.MustAddress("0x3333333333333333333333333333333333333333"),
		ExpirySec:        1900000000,
		Scope:            "read_only",
		Label:            "scoped",
	}
	got, err := api.RegisterScopedSessionKey(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, "txn-42", got.TransactionID)
	assert.Equal(t, "read_only", got.Scope)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "read_only", params["scope"])
	assert.Equal(t, "scoped", params["label"])
	assert.Equal(t, float64(1900000000), params["expiry_sec"])
	_, hasSignedTx := params["signed_raw_tx"]
	assert.False(t, hasSignedTx, "signed_raw_tx should be omitted for non-admin scope")
}

func TestRegisterScopedSessionKey_AdminCarriesSignedTx(t *testing.T) {
	api, ft := newAPI(t, true, 0)
	ft.HandleResult("private/register_scoped_session_key", map[string]any{
		"public_session_key": "0x3333333333333333333333333333333333333333",
		"scope":              "admin",
		"expiry_sec":         int64(1900000000),
		"ip_whitelist":       []any{},
		"transaction_id":     "txn-43",
		"label":              "",
	})
	in := types.RegisterScopedSessionKeyInput{
		PublicSessionKey: types.MustAddress("0x3333333333333333333333333333333333333333"),
		ExpirySec:        1900000000,
		Scope:            "admin",
		SignedRawTx:      "0xdeadbeef",
	}
	_, err := api.RegisterScopedSessionKey(context.Background(), in)
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "0xdeadbeef", params["signed_raw_tx"])
}

func TestRegisterScopedSessionKey_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.RegisterScopedSessionKey(context.Background(), types.RegisterScopedSessionKeyInput{})
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
}
