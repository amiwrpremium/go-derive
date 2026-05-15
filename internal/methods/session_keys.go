// Package methods is the shared implementation of every JSON-RPC method
// Derive exposes. Both pkg/rest.Client and pkg/ws.Client embed *API so that
// each method is defined exactly once, parameterised by the underlying
// transport.
package methods

import (
	"context"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// SessionKeys lists every session key registered for the configured
// signer's owner wallet. Private. Wraps `private/session_keys`.
//
// The response includes unactivated and expired keys — filter on
// [types.SessionKey.ExpirySec] at the call site if only live keys are
// of interest.
//
// Minimum session-key permission level: read_only.
func (a *API) SessionKeys(ctx context.Context) ([]types.SessionKey, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	params := map[string]any{"wallet": a.Signer.OwnerAddress().Hex()}
	var resp struct {
		PublicSessionKeys []types.SessionKey `json:"public_session_keys"`
	}
	if err := a.call(ctx, "private/session_keys", params, &resp); err != nil {
		return nil, err
	}
	return resp.PublicSessionKeys, nil
}

// EditSessionKey updates the label, IP allow-list, or disabled flag
// on an existing session key. Private. Wraps `private/edit_session_key`.
//
// Only fields explicitly set on [types.EditSessionKeyInput] are sent;
// leaving a pointer field nil omits that property and preserves its
// server-side value. Pass an empty IP allow-list (not nil) to clear
// the allow-list.
//
// Admin-scope keys cannot be disabled here; use
// `public/deregister_session_key` instead.
//
// Minimum session-key permission level: admin.
func (a *API) EditSessionKey(ctx context.Context, in types.EditSessionKeyInput) (types.SessionKey, error) {
	if err := a.requireSigner(); err != nil {
		return types.SessionKey{}, err
	}
	params := map[string]any{
		"wallet":             a.Signer.OwnerAddress().Hex(),
		"public_session_key": in.PublicSessionKey.String(),
	}
	if in.Disable {
		params["disable"] = true
	}
	if in.Label != nil {
		params["label"] = *in.Label
	}
	if in.IPWhitelist != nil {
		params["ip_whitelist"] = *in.IPWhitelist
	}
	var resp types.SessionKey
	if err := a.call(ctx, "private/edit_session_key", params, &resp); err != nil {
		return types.SessionKey{}, err
	}
	return resp, nil
}

// RegisterScopedSessionKey registers a new scoped session key for the
// configured signer's owner wallet. Private. Wraps
// `private/register_scoped_session_key`.
//
// The on-chain registration is asynchronous. The returned
// TransactionID can be polled via [API.GetTransaction]. The session
// key is not usable until that transaction settles.
//
// For ADMIN-scope keys, [types.RegisterScopedSessionKeyInput.SignedRawTx]
// is required and must carry a pre-signed RLP-encoded Ethereum
// transaction (matching Python's
// `w3.eth.account.sign_transaction(tx, priv_key).rawTransaction.hex()`).
//
// Minimum session-key permission level: admin.
func (a *API) RegisterScopedSessionKey(ctx context.Context, in types.RegisterScopedSessionKeyInput) (types.RegisterScopedSessionKeyResult, error) {
	if err := a.requireSigner(); err != nil {
		return types.RegisterScopedSessionKeyResult{}, err
	}
	params := map[string]any{
		"wallet":             a.Signer.OwnerAddress().Hex(),
		"public_session_key": in.PublicSessionKey.String(),
		"expiry_sec":         in.ExpirySec,
	}
	if in.Scope != "" {
		params["scope"] = in.Scope
	}
	if in.Label != "" {
		params["label"] = in.Label
	}
	if in.IPWhitelist != nil {
		params["ip_whitelist"] = in.IPWhitelist
	}
	if in.SignedRawTx != "" {
		params["signed_raw_tx"] = in.SignedRawTx
	}
	var resp types.RegisterScopedSessionKeyResult
	if err := a.call(ctx, "private/register_scoped_session_key", params, &resp); err != nil {
		return types.RegisterScopedSessionKeyResult{}, err
	}
	return resp, nil
}
