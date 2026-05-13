// Package types declares the domain types used in REST and WebSocket
// requests and responses.
package types

// SessionKey is one entry in the response from `private/session_keys` and
// `private/edit_session_key`. It describes a session key registered for a
// wallet, including its scope, IP allow-list, label, and lifecycle
// timestamps.
type SessionKey struct {
	// PublicSessionKey is the session key's EOA address.
	PublicSessionKey Address `json:"public_session_key"`
	// Label is the user-defined session-key label.
	Label string `json:"label"`
	// Scope is the permission level granted to the session key
	// ("admin", "account", or "read_only").
	Scope string `json:"scope"`
	// ExpirySec is the session-key expiry as a Unix timestamp in
	// seconds.
	ExpirySec int64 `json:"expiry_sec"`
	// RegisteredSec is when the session key was registered, as a
	// Unix timestamp in seconds.
	RegisteredSec int64 `json:"registered_sec"`
	// IPWhitelist is the optional list of source IPs allowed to use
	// this session key. An empty list means any IP is allowed.
	IPWhitelist []string `json:"ip_whitelist"`
}

// EditSessionKeyInput is the input to `private/edit_session_key`. Only
// PublicSessionKey is required; the rest are nullable on the wire and
// optional here — leave a field zero to skip changing it.
type EditSessionKeyInput struct {
	// PublicSessionKey is the session key (EOA address) to edit.
	PublicSessionKey Address
	// Disable disables the key when true. Only allowed for non-admin
	// keys; admin keys must be deregistered via
	// `public/deregister_session_key`.
	Disable bool
	// Label, when non-nil, replaces the existing label.
	Label *string
	// IPWhitelist, when non-nil, replaces the existing IP allow-list.
	// Pass an empty slice to clear the allow-list.
	IPWhitelist *[]string
}

// RegisterScopedSessionKeyInput is the input to
// `private/register_scoped_session_key`. The on-chain registration is
// asynchronous: the response carries a transaction id that can be
// polled via `public/get_transaction`.
type RegisterScopedSessionKeyInput struct {
	// PublicSessionKey is the EOA address to register.
	PublicSessionKey Address
	// ExpirySec is the Unix-seconds expiry of the session key.
	ExpirySec int64
	// Scope is the permission level requested. Defaults server-side
	// to "read_only" if empty.
	Scope string
	// Label is the optional user-defined session-key label.
	Label string
	// IPWhitelist is the optional list of source IPs. Empty means
	// any IP is allowed.
	IPWhitelist []string
	// SignedRawTx is a signed RLP-encoded ETH transaction (hex
	// string) required for ADMIN-scope registrations and ignored
	// otherwise.
	SignedRawTx string
}

// RegisterScopedSessionKeyResult is the response from
// `private/register_scoped_session_key`. The session key is not
// usable until the carried transaction settles on chain — poll the
// returned TransactionID via `public/get_transaction`.
type RegisterScopedSessionKeyResult struct {
	// PublicSessionKey is the registered session key.
	PublicSessionKey Address `json:"public_session_key"`
	// Label is the session-key label.
	Label string `json:"label,omitempty"`
	// Scope is the permission level granted.
	Scope string `json:"scope"`
	// ExpirySec is the session-key expiry, Unix seconds.
	ExpirySec int64 `json:"expiry_sec"`
	// IPWhitelist is the optional IP allow-list. Empty means any IP.
	IPWhitelist []string `json:"ip_whitelist,omitempty"`
	// TransactionID is the engine-side transaction id for the
	// on-chain registration. Poll via `public/get_transaction`.
	TransactionID string `json:"transaction_id"`
}
