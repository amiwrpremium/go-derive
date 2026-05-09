// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shapes for `private/replace`,
// `private/order_debug`, and `private/cancel_by_nonce`.
package types

import "encoding/json"

// ReplaceResult is the response of `private/replace`. The endpoint
// cancels one outstanding order and submits a replacement atomically;
// the response carries both the cancelled and the (optional) newly
// created order, plus any trades the new order matched, plus an
// optional engine error if the new order was rejected.
//
// The shape mirrors `PrivateReplaceResultSchema` in Derive's v2.2
// OpenAPI spec — `cancelled_order` and `trades` are required (the
// trades array is empty when the replacement didn't fill); `order`
// and `create_order_error` are mutually-exclusive optionals.
type ReplaceResult struct {
	// CancelledOrder is the order that was dropped.
	CancelledOrder Order `json:"cancelled_order"`
	// Order is the replacement order. Nil when CreateOrderError is
	// non-nil — the engine cancelled the old order but rejected the
	// new one.
	Order *Order `json:"order,omitempty"`
	// CreateOrderError is the engine error if the replacement order
	// was rejected. Non-nil only on failure.
	CreateOrderError *RPCError `json:"create_order_error,omitempty"`
	// Trades is the list of fills the new order produced before any
	// remaining size came to rest. Always present per the OAS — an
	// empty array when the replacement didn't fill (no `omitempty`
	// because the wire field is required).
	Trades []Trade `json:"trades"`
}

// RPCError is the JSON-RPC error envelope embedded in
// `ReplaceResult.CreateOrderError`. It mirrors `RPCErrorFormatSchema`
// in the OAS.
type RPCError struct {
	// Code is the JSON-RPC error code.
	Code int `json:"code"`
	// Message is the engine's short description.
	Message string `json:"message"`
	// Data is the optional structured payload. The wire field is
	// nullable; absent → empty string.
	Data string `json:"data,omitempty"`
}

// OrderDebugResult is the response of `private/order_debug`. The
// endpoint mirrors `private/order` but returns the engine's internal
// hashing artefacts instead of placing the order — useful for
// validating signatures in CI.
//
// The shape mirrors `PrivateOrderDebugResultSchema` in Derive's v2.2
// OpenAPI spec.
type OrderDebugResult struct {
	// ActionHash is the keccak hash of the action data.
	ActionHash string `json:"action_hash"`
	// EncodedData is the ABI-encoded order data.
	EncodedData string `json:"encoded_data"`
	// EncodedDataHashed is the keccak hash of EncodedData.
	EncodedDataHashed string `json:"encoded_data_hashed"`
	// TypedDataHash is the EIP-712 typed-data hash the engine
	// computed.
	TypedDataHash string `json:"typed_data_hash"`
	// RawData is the engine's internal view of the signed order.
	RawData OrderDebugRawData `json:"raw_data"`
}

// OrderDebugRawData is the engine's internal `SignedTradeOrderSchema`
// view returned by `private/order_debug`.
//
// `Data` is the per-module payload. For order-flow debug calls
// (which this method always is) it's the engine's
// `TradeModuleDataSchema` shape:
//
//	{
//	    asset:           string,
//	    desired_amount:  decimal,
//	    is_bid:          bool,
//	    limit_price:     decimal,
//	    recipient_id:    int,
//	    sub_id:          int,
//	    trade_id:        string,
//	    worst_fee:       decimal,
//	}
//
// We keep it as a raw payload because:
//
//   - The wire field names differ from the [pkg/auth.TradeModuleData]
//     input shape ([pkg/auth] uses `Amount`/`MaxFee`; the engine
//     emits `desired_amount`/`worst_fee`), so a typed Go shape
//     wouldn't reuse the existing input type cleanly.
//   - Callers using `private/order_debug` are typically
//     pre-flighting orders and only care that the call succeeded;
//     the bytes are surfaced for diagnostic logging if needed.
//
// Decode against the schema above when you need the structured
// payload.
type OrderDebugRawData struct {
	// Data is the engine's per-module signed payload. See the type
	// comment above for the documented `TradeModuleDataSchema`
	// shape that this method's response carries.
	Data json.RawMessage `json:"data"`
	// Expiry is the action's expiry (Unix seconds).
	Expiry int64 `json:"expiry"`
	// IsAtomicSigning reports whether the engine treated the
	// signature as atomic.
	IsAtomicSigning bool `json:"is_atomic_signing"`
	// Module is the on-chain module address as a string.
	Module string `json:"module"`
	// Nonce is the action nonce.
	Nonce int64 `json:"nonce"`
	// Owner is the smart-account owner address.
	Owner string `json:"owner"`
	// Signature is the EIP-712 signature.
	Signature string `json:"signature"`
	// Signer is the address that signed.
	Signer string `json:"signer"`
	// SubaccountID is the placing subaccount.
	SubaccountID int64 `json:"subaccount_id"`
}

// CancelByNonceResult is the response of `private/cancel_by_nonce`.
// It reports how many orders matched the (instrument, nonce) tuple
// and were cancelled.
type CancelByNonceResult struct {
	// CancelledOrders is the count of orders cancelled.
	CancelledOrders int64 `json:"cancelled_orders"`
}
