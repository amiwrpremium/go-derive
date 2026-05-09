package contracts_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/contracts"
)

type stubDepositor struct{}

func (stubDepositor) Deposit(context.Context, derive.Address, decimal.Decimal, int64) (derive.TxHash, error) {
	return derive.TxHash{}, contracts.ErrNotImplemented
}

func TestDepositor_InterfaceSatisfied(t *testing.T) {
	var d contracts.Depositor = stubDepositor{}
	tx, err := d.Deposit(context.Background(), derive.Address{}, decimal.Zero, 0)
	assert.ErrorIs(t, err, contracts.ErrNotImplemented)
	assert.True(t, tx.IsZero())
}

func TestErrNotImplemented_Message(t *testing.T) {
	assert.NotNil(t, contracts.ErrNotImplemented)
	assert.Contains(t, contracts.ErrNotImplemented.Error(), "not implemented")
}
