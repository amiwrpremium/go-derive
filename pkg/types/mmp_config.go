// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the input DTO for `private/set_mmp_config` —
// previously declared in internal/methods, lifted here so callers
// only need to import pkg/types for the SDK's domain types.
package types

// MMPConfig is the input to `private/set_mmp_config` — Market Maker
// Protection rules for a single (subaccount, currency) pair.
type MMPConfig struct {
	Currency        string `json:"currency"`
	MMPFrozenTimeMs int64  `json:"mmp_frozen_time"`
	MMPIntervalMs   int64  `json:"mmp_interval"`
	MMPAmountLimit  string `json:"mmp_amount_limit,omitempty"`
	MMPDeltaLimit   string `json:"mmp_delta_limit,omitempty"`
}

// Validate performs schema-level checks on the receiver. Returns nil
// on success or an error wrapping [ErrInvalidParams]. The two limit
// fields are decimal strings on the wire and remain unparsed here.
func (c MMPConfig) Validate() error {
	if c.Currency == "" {
		return invalidParam("currency", "required")
	}
	if c.MMPFrozenTimeMs < 0 {
		return invalidParam("mmp_frozen_time", "must be non-negative")
	}
	if c.MMPIntervalMs < 0 {
		return invalidParam("mmp_interval", "must be non-negative")
	}
	return nil
}
