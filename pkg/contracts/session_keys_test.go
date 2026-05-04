package contracts_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/contracts"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

type stubSessionKeyManager struct{}

func (stubSessionKeyManager) Register(context.Context, types.Address, time.Time) (types.TxHash, error) {
	return types.TxHash{}, contracts.ErrNotImplemented
}
func (stubSessionKeyManager) Revoke(context.Context, types.Address) (types.TxHash, error) {
	return types.TxHash{}, contracts.ErrNotImplemented
}

func TestSessionKeyManager_RegisterReturnsErr(t *testing.T) {
	var m contracts.SessionKeyManager = stubSessionKeyManager{}
	_, err := m.Register(context.Background(), types.Address{}, time.Now())
	assert.ErrorIs(t, err, contracts.ErrNotImplemented)
}

func TestSessionKeyManager_RevokeReturnsErr(t *testing.T) {
	var m contracts.SessionKeyManager = stubSessionKeyManager{}
	_, err := m.Revoke(context.Background(), types.Address{})
	assert.ErrorIs(t, err, contracts.ErrNotImplemented)
}
