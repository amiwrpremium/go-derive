package contracts_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/contracts"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

type stubDepositor struct{}

func (stubDepositor) Deposit(context.Context, types.Address, decimal.Decimal, int64) (types.TxHash, error) {
	return types.TxHash{}, contracts.ErrNotImplemented
}

func TestDepositor_InterfaceSatisfied(t *testing.T) {
	var d contracts.Depositor = stubDepositor{}
	tx, err := d.Deposit(context.Background(), types.Address{}, decimal.Zero, 0)
	assert.ErrorIs(t, err, contracts.ErrNotImplemented)
	assert.True(t, tx.IsZero())
}

func TestErrNotImplemented_Message(t *testing.T) {
	assert.NotNil(t, contracts.ErrNotImplemented)
	assert.Contains(t, contracts.ErrNotImplemented.Error(), "not implemented")
}
