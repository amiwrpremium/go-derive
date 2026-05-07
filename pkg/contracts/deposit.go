// Package contracts hosts on-chain helper interfaces — deposits,
// withdrawals, and session-key lifecycle — for Derive's smart-account
// model.
//
// # Status
//
// The package is intentionally a stub: the JSON-RPC layer
// ([github.com/amiwrpremium/go-derive/pkg/rest] and
// [github.com/amiwrpremium/go-derive/pkg/ws]) is sufficient to trade once
// collateral has been deposited via the Derive UI or another EVM tool.
// Every interface in this package is declared so that consumers can write
// code against the API today against a stable shape.
//
// All methods return [ErrNotImplemented].
package contracts

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// Depositor is the interface for crediting collateral into a subaccount.
type Depositor interface {
	// Deposit credits collateral into a subaccount on Derive Chain. It
	// submits an ERC-20 approve+deposit pair as one logical operation and
	// returns the deposit transaction hash on success.
	//
	// Returns [ErrNotImplemented].
	Deposit(ctx context.Context, asset types.Address, amount decimal.Decimal, subaccount int64) (types.TxHash, error)
}

// ErrNotImplemented is returned by every method on the stubs in this
// package. Use it as a sentinel via errors.Is:
//
//	if errors.Is(err, contracts.ErrNotImplemented) { ... }
var ErrNotImplemented = errors.New("contracts: on-chain helpers are not implemented; see pkg/contracts/doc.go")
