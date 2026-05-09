package contracts_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/contracts"
)

type stubWithdrawer struct{}

func (stubWithdrawer) Withdraw(context.Context, derive.Address, decimal.Decimal, int64) (derive.TxHash, error) {
	return derive.TxHash{}, contracts.ErrNotImplemented
}

func TestWithdrawer_InterfaceSatisfied(t *testing.T) {
	var w contracts.Withdrawer = stubWithdrawer{}
	_, err := w.Withdraw(context.Background(), derive.Address{}, decimal.Zero, 0)
	assert.ErrorIs(t, err, contracts.ErrNotImplemented)
}

func TestWithdrawer_NonZeroArgs(t *testing.T) {
	var w contracts.Withdrawer = stubWithdrawer{}
	addr := derive.MustAddress("0x1111111111111111111111111111111111111111")
	tx, err := w.Withdraw(context.Background(), addr, decimal.RequireFromString("1.5"), 42)
	assert.ErrorIs(t, err, contracts.ErrNotImplemented)
	assert.True(t, tx.IsZero())
}
