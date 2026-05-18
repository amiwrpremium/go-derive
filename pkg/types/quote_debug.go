// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shapes for the RFQ debug endpoints
// `public/send_quote_debug` and `public/execute_quote_debug`.
package types

// SendQuoteDebugResult is the response of `public/send_quote_debug`.
// The endpoint mirrors `private/send_quote` but returns the engine's
// internal hashing artefacts instead of placing the quote — useful
// for validating per-quote EIP-712 signatures in CI before paying
// gas or risking a live RFQ-state desync.
//
// The shape mirrors `PublicSendQuoteDebugResultSchema` in Derive's
// v2.2 OpenAPI spec.
type SendQuoteDebugResult struct {
	// ActionHash is the keccak hash of the action data.
	ActionHash string `json:"action_hash"`
	// EncodedData is the ABI-encoded quote data.
	EncodedData string `json:"encoded_data"`
	// EncodedDataHashed is the keccak hash of EncodedData.
	EncodedDataHashed string `json:"encoded_data_hashed"`
	// TypedDataHash is the EIP-712 typed-data hash the engine
	// computed.
	TypedDataHash string `json:"typed_data_hash"`
}

// ExecuteQuoteDebugResult is the response of
// `public/execute_quote_debug`. Same four hashes as
// [SendQuoteDebugResult], plus the ABI-encoded leg payload:
// `execute_quote` inverts each leg's direction relative to
// `send_quote`, so the encoded leg bytes are materially different
// and surfaced separately for verification.
//
// The shape mirrors `PublicExecuteQuoteDebugResultSchema` in
// Derive's v2.2 OpenAPI spec.
type ExecuteQuoteDebugResult struct {
	// ActionHash is the keccak hash of the action data.
	ActionHash string `json:"action_hash"`
	// EncodedData is the ABI-encoded execute-quote data.
	EncodedData string `json:"encoded_data"`
	// EncodedDataHashed is the keccak hash of EncodedData.
	EncodedDataHashed string `json:"encoded_data_hashed"`
	// EncodedLegs is the ABI-encoded leg payload — unique to
	// `execute_quote_debug` because executes invert each leg's
	// direction.
	EncodedLegs string `json:"encoded_legs"`
	// LegsHash is the keccak hash of EncodedLegs.
	LegsHash string `json:"legs_hash"`
	// TypedDataHash is the EIP-712 typed-data hash the engine
	// computed.
	TypedDataHash string `json:"typed_data_hash"`
}
