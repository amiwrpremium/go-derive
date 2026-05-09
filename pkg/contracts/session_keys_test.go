package contracts_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/contracts"
)

type stubSessionKeyManager struct{}

func (stubSessionKeyManager) Register(context.Context, derive.Address, time.Time) (derive.TxHash, error) {
	return derive.TxHash{}, contracts.ErrNotImplemented
}
func (stubSessionKeyManager) Revoke(context.Context, derive.Address) (derive.TxHash, error) {
	return derive.TxHash{}, contracts.ErrNotImplemented
}

func TestSessionKeyManager_RegisterReturnsErr(t *testing.T) {
	var m contracts.SessionKeyManager = stubSessionKeyManager{}
	_, err := m.Register(context.Background(), derive.Address{}, time.Now())
	assert.ErrorIs(t, err, contracts.ErrNotImplemented)
}

func TestSessionKeyManager_RevokeReturnsErr(t *testing.T) {
	var m contracts.SessionKeyManager = stubSessionKeyManager{}
	_, err := m.Revoke(context.Background(), derive.Address{})
	assert.ErrorIs(t, err, contracts.ErrNotImplemented)
}
