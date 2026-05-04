package auth_test

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

var goodAsset = common.HexToAddress("0x1111111111111111111111111111111111111111")

func validTrade() auth.TradeModuleData {
	return auth.TradeModuleData{
		Asset:       goodAsset,
		SubID:       7,
		LimitPrice:  decimal.RequireFromString("100"),
		Amount:      decimal.RequireFromString("0.5"),
		MaxFee:      decimal.RequireFromString("1"),
		RecipientID: 42,
		IsBid:       true,
	}
}

func TestTradeModuleData_Validate_Happy(t *testing.T) {
	require.NoError(t, validTrade().Validate())
}

func TestTradeModuleData_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*auth.TradeModuleData)
		want string
	}{
		{"zero asset", func(t *auth.TradeModuleData) { t.Asset = common.Address{} }, "asset"},
		{"zero price", func(t *auth.TradeModuleData) { t.LimitPrice = decimal.Zero }, "limit_price"},
		{"negative price", func(t *auth.TradeModuleData) { t.LimitPrice = decimal.RequireFromString("-1") }, "limit_price"},
		{"zero amount", func(t *auth.TradeModuleData) { t.Amount = decimal.Zero }, "amount"},
		{"negative fee", func(t *auth.TradeModuleData) { t.MaxFee = decimal.RequireFromString("-1") }, "max_fee"},
		{"negative recipient", func(t *auth.TradeModuleData) { t.RecipientID = -1 }, "recipient_id"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := validTrade()
			c.mut(&d)
			err := d.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, auth.ErrInvalidInput))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func validTransfer() auth.TransferModuleData {
	return auth.TransferModuleData{
		ToSubaccount: 1,
		Asset:        goodAsset,
		SubID:        3,
		Amount:       decimal.RequireFromString("10"),
	}
}

func TestTransferModuleData_Validate_Happy(t *testing.T) {
	require.NoError(t, validTransfer().Validate())
}

func TestTransferModuleData_Validate_AllowsNegativeAmount(t *testing.T) {
	d := validTransfer()
	d.Amount = decimal.RequireFromString("-5")
	require.NoError(t, d.Validate())
}

func TestTransferModuleData_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*auth.TransferModuleData)
		want string
	}{
		{"zero asset", func(t *auth.TransferModuleData) { t.Asset = common.Address{} }, "asset"},
		{"negative subaccount", func(t *auth.TransferModuleData) { t.ToSubaccount = -1 }, "to_subaccount"},
		{"zero amount", func(t *auth.TransferModuleData) { t.Amount = decimal.Zero }, "amount"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := validTransfer()
			c.mut(&d)
			err := d.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, auth.ErrInvalidInput))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func validAction() auth.ActionData {
	return auth.ActionData{
		SubaccountID: 7,
		Nonce:        42,
		Module:       goodAsset,
		Expiry:       1_700_000_000,
		Owner:        common.HexToAddress("0x2222222222222222222222222222222222222222"),
		Signer:       common.HexToAddress("0x3333333333333333333333333333333333333333"),
	}
}

func TestActionData_Validate_Happy(t *testing.T) {
	require.NoError(t, validAction().Validate())
}

func TestActionData_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*auth.ActionData)
		want string
	}{
		{"negative subaccount", func(a *auth.ActionData) { a.SubaccountID = -1 }, "subaccount_id"},
		{"zero module", func(a *auth.ActionData) { a.Module = common.Address{} }, "module"},
		{"zero owner", func(a *auth.ActionData) { a.Owner = common.Address{} }, "owner"},
		{"zero signer", func(a *auth.ActionData) { a.Signer = common.Address{} }, "signer"},
		{"zero expiry", func(a *auth.ActionData) { a.Expiry = 0 }, "expiry"},
		{"negative expiry", func(a *auth.ActionData) { a.Expiry = -1 }, "expiry"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := validAction()
			c.mut(&d)
			err := d.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, auth.ErrInvalidInput))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}
