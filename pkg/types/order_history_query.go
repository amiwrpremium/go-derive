// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the query DTO for `private/get_order_history` —
// previously declared in internal/methods, lifted here so callers
// only need to import pkg/types for the SDK's domain types.
package types

// OrderHistoryQuery narrows a paginated `private/get_order_history`
// request. FromTimestamp / ToTimestamp form a closed window in
// milliseconds since the Unix epoch; either side can be zero to defer
// to the server-side default (0 / current time). Wallet, when
// non-empty, queries across every subaccount under that wallet —
// when empty, the configured subaccount is used.
type OrderHistoryQuery struct {
	FromTimestamp MillisTime
	ToTimestamp   MillisTime
	Wallet        string
}
