package derive_test

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive"
)

type stubDepositor struct{}

func (stubDepositor) Deposit(context.Context, derive.Address, decimal.Decimal, int64) (derive.TxHash, error) {
	return derive.TxHash{}, derive.ErrNotImplemented
}

func TestDepositor_InterfaceSatisfied(t *testing.T) {
	var d derive.Depositor = stubDepositor{}
	tx, err := d.Deposit(context.Background(), derive.Address{}, decimal.Zero, 0)
	assert.ErrorIs(t, err, derive.ErrNotImplemented)
	assert.True(t, tx.IsZero())
}

func TestErrNotImplemented_Message(t *testing.T) {
	assert.NotNil(t, derive.ErrNotImplemented)
	assert.Contains(t, derive.ErrNotImplemented.Error(), "not implemented")
}

type stubWithdrawer struct{}

func (stubWithdrawer) Withdraw(context.Context, derive.Address, decimal.Decimal, int64) (derive.TxHash, error) {
	return derive.TxHash{}, derive.ErrNotImplemented
}

func TestWithdrawer_InterfaceSatisfied(t *testing.T) {
	var w derive.Withdrawer = stubWithdrawer{}
	_, err := w.Withdraw(context.Background(), derive.Address{}, decimal.Zero, 0)
	assert.ErrorIs(t, err, derive.ErrNotImplemented)
}

func TestWithdrawer_NonZeroArgs(t *testing.T) {
	var w derive.Withdrawer = stubWithdrawer{}
	addr := derive.MustAddress("0x1111111111111111111111111111111111111111")
	tx, err := w.Withdraw(context.Background(), addr, decimal.RequireFromString("1.5"), 42)
	assert.ErrorIs(t, err, derive.ErrNotImplemented)
	assert.True(t, tx.IsZero())
}

type stubSessionKeyManager struct{}

func (stubSessionKeyManager) Register(context.Context, derive.Address, time.Time) (derive.TxHash, error) {
	return derive.TxHash{}, derive.ErrNotImplemented
}
func (stubSessionKeyManager) Revoke(context.Context, derive.Address) (derive.TxHash, error) {
	return derive.TxHash{}, derive.ErrNotImplemented
}

func TestSessionKeyManager_RegisterReturnsErr(t *testing.T) {
	var m derive.SessionKeyManager = stubSessionKeyManager{}
	_, err := m.Register(context.Background(), derive.Address{}, time.Now())
	assert.ErrorIs(t, err, derive.ErrNotImplemented)
}

func TestSessionKeyManager_RevokeReturnsErr(t *testing.T) {
	var m derive.SessionKeyManager = stubSessionKeyManager{}
	_, err := m.Revoke(context.Background(), derive.Address{})
	assert.ErrorIs(t, err, derive.ErrNotImplemented)
}
