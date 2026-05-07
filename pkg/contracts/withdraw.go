// Package contracts.
package contracts

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// Withdrawer is the interface for withdrawing collateral from a subaccount.
type Withdrawer interface {
	// Withdraw debits collateral from a subaccount and queues the on-chain
	// transfer back to the owner wallet. It returns the withdrawal
	// transaction hash on success.
	//
	// Returns [ErrNotImplemented].
	Withdraw(ctx context.Context, asset types.Address, amount decimal.Decimal, subaccount int64) (types.TxHash, error)
}
