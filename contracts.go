// On-chain helper interfaces — deposits, withdrawals, and session-key
// lifecycle — for Derive's smart-account model live in this file.
//
// # Status
//
// These interfaces are intentional stubs: the JSON-RPC layer ([RestClient]
// / [WsClient]) is sufficient to trade once collateral has been deposited
// via the Derive UI or another EVM tool. The shapes are declared so that
// consumers can write code against the API today against a stable
// interface.
//
// All methods return [ErrNotImplemented].

package derive

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// ErrNotImplemented is returned by every method on the on-chain helper
// stubs in this package. Use it as a sentinel via errors.Is:
//
//	if errors.Is(err, derive.ErrNotImplemented) { ... }
var ErrNotImplemented = errors.New("derive: on-chain helpers are not implemented; see contracts.go")

// Depositor is the interface for crediting collateral into a subaccount.
type Depositor interface {
	// Deposit credits collateral into a subaccount on Derive Chain. It
	// submits an ERC-20 approve+deposit pair as one logical operation and
	// returns the deposit transaction hash on success.
	//
	// Returns [ErrNotImplemented].
	Deposit(ctx context.Context, asset Address, amount decimal.Decimal, subaccount int64) (TxHash, error)
}

// Withdrawer is the interface for withdrawing collateral from a subaccount.
type Withdrawer interface {
	// Withdraw debits collateral from a subaccount and queues the on-chain
	// transfer back to the owner wallet. It returns the withdrawal
	// transaction hash on success.
	//
	// Returns [ErrNotImplemented].
	Withdraw(ctx context.Context, asset Address, amount decimal.Decimal, subaccount int64) (TxHash, error)
}

// SessionKeyManager is the contract for the session-key lifecycle.
// Session keys are addresses authorised to sign Derive actions on behalf
// of the owner wallet; they limit blast radius if a hot key is
// compromised.
type SessionKeyManager interface {
	// Register adds a session key authorised to sign actions on behalf of
	// the owner wallet, valid until expiry. It returns the registration
	// transaction hash on success.
	//
	// Returns [ErrNotImplemented].
	Register(ctx context.Context, sessionKey Address, expiry time.Time) (TxHash, error)

	// Revoke immediately deauthorises a session key. It returns the
	// revocation transaction hash on success.
	//
	// Returns [ErrNotImplemented].
	Revoke(ctx context.Context, sessionKey Address) (TxHash, error)
}
