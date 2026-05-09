// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// All numeric fields use [Decimal], a thin wrapper around shopspring/decimal,
// so price/size/fee values never lose precision through float64 round-trips.
// On the wire, [Decimal] reads and writes JSON strings (Derive's preferred
// representation); a fallback path also accepts JSON numbers for resilience.
//
// Identifier types ([Address], [TxHash], [MillisTime]) carry the same
// round-trip guarantees: each one preserves the canonical wire format
// regardless of how Go marshals the surrounding struct.
//
// # Why named types
//
// Plain string and int64 fields would parse just fine, but named types let
// the SDK enforce invariants at construction time (NewAddress checksum
// check, NewDecimal precision check) and let callers tell at a glance which
// values are amounts vs prices vs subaccount ids.
package types

// DepositTx records a single deposit into a subaccount.
//
// Returned by `private/get_deposit_history`; also delivered on the
// account-balance channel as deposits finalize.
type DepositTx struct {
	// TxHash is the on-chain deposit transaction hash.
	TxHash TxHash `json:"tx_hash"`
	// Asset is the deposited asset's symbol (e.g. "USDC").
	Asset string `json:"asset"`
	// Amount is the deposited quantity.
	Amount Decimal `json:"amount"`
	// SubaccountID is the credited subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Status is the lifecycle state ("pending", "completed", etc.).
	Status string `json:"status"`
	// Timestamp is when the deposit was first observed.
	Timestamp MillisTime `json:"timestamp"`
}

// Transaction is the response of `public/get_transaction`. It carries
// the engine's view of one Derive transaction by its server-side id.
//
// Mirrors `PublicGetTransactionResultSchema` in Derive's v2.2
// OpenAPI spec.
type Transaction struct {
	// Data is the request body that produced the transaction.
	// Surfaced as a string by the API even though it's a structured
	// payload — decode at the call site if needed.
	Data string `json:"data"`
	// ErrorLog is populated if the transaction failed. Nullable on
	// the wire; absent decodes to empty string.
	ErrorLog string `json:"error_log,omitempty"`
	// Status is the transaction lifecycle state: "requested",
	// "pending", "settled", "reverted", "ignored", or "timed_out".
	Status string `json:"status"`
	// TransactionHash is the on-chain transaction hash. Nullable on
	// the wire (the server returns it once submitted); absent
	// decodes to empty string.
	TransactionHash string `json:"transaction_hash,omitempty"`
}

// WithdrawTx records a single withdrawal from a subaccount.
//
// Withdrawals are two-phase: first the subaccount is debited (status
// pending), then the on-chain transfer is dispatched (status completed).
type WithdrawTx struct {
	// TxHash is the on-chain withdrawal transaction hash.
	TxHash TxHash `json:"tx_hash"`
	// Asset is the withdrawn asset's symbol.
	Asset string `json:"asset"`
	// Amount is the withdrawn quantity.
	Amount Decimal `json:"amount"`
	// SubaccountID is the debited subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Status is the lifecycle state.
	Status string `json:"status"`
	// Timestamp is when the withdrawal was first observed.
	Timestamp MillisTime `json:"timestamp"`
}
