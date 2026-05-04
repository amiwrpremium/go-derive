package contracts_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/contracts"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

type stubWithdrawer struct{}

func (stubWithdrawer) Withdraw(context.Context, types.Address, decimal.Decimal, int64) (types.TxHash, error) {
	return types.TxHash{}, contracts.ErrNotImplemented
}

func TestWithdrawer_InterfaceSatisfied(t *testing.T) {
	var w contracts.Withdrawer = stubWithdrawer{}
	_, err := w.Withdraw(context.Background(), types.Address{}, decimal.Zero, 0)
	assert.ErrorIs(t, err, contracts.ErrNotImplemented)
}

func TestWithdrawer_NonZeroArgs(t *testing.T) {
	var w contracts.Withdrawer = stubWithdrawer{}
	addr := types.MustAddress("0x1111111111111111111111111111111111111111")
	tx, err := w.Withdraw(context.Background(), addr, decimal.RequireFromString("1.5"), 42)
	assert.ErrorIs(t, err, contracts.ErrNotImplemented)
	assert.True(t, tx.IsZero())
}
