//go:build integration

package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestPrivate_GetSubaccount(t *testing.T) {
	env := loadEnv(t)
	c := newAuthRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	sa, err := c.GetSubaccount(ctx)
	require.NoError(t, err)
	assert.Equal(t, env.subaccount, sa.SubaccountID)
}

func TestPrivate_GetSubaccounts(t *testing.T) {
	env := loadEnv(t)
	c := newAuthRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	ids, err := c.GetSubaccounts(ctx)
	require.NoError(t, err)
	assert.Contains(t, ids, env.subaccount, "configured subaccount should be in the wallet's list")
}

func TestPrivate_GetOpenOrders(t *testing.T) {
	env := loadEnv(t)
	c := newAuthRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	got, err := c.GetOpenOrders(ctx)
	require.NoError(t, err)
	for _, o := range got {
		assert.Equal(t, env.subaccount, o.SubaccountID)
	}
}

func TestPrivate_GetPositions(t *testing.T) {
	env := loadEnv(t)
	c := newAuthRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	_, err := c.GetPositions(ctx)
	require.NoError(t, err)
}

func TestPrivate_GetCollateral(t *testing.T) {
	env := loadEnv(t)
	c := newAuthRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	_, err := c.GetCollateral(ctx)
	require.NoError(t, err)
}

func TestPrivate_GetTradeHistory(t *testing.T) {
	env := loadEnv(t)
	c := newAuthRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	_, page, err := c.GetTradeHistory(ctx, types.PageRequest{PageSize: 10})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, page.Count, 0)
}

func TestPrivate_GetDepositHistory(t *testing.T) {
	env := loadEnv(t)
	c := newAuthRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	_, _, err := c.GetDepositHistory(ctx, types.PageRequest{PageSize: 10})
	require.NoError(t, err)
}

func TestPrivate_GetWithdrawalHistory(t *testing.T) {
	env := loadEnv(t)
	c := newAuthRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	_, _, err := c.GetWithdrawalHistory(ctx, types.PageRequest{PageSize: 10})
	require.NoError(t, err)
}
