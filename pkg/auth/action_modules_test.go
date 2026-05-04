package auth_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func baseTrade() auth.TradeModuleData {
	return auth.TradeModuleData{
		Asset:       common.HexToAddress("0x1111111111111111111111111111111111111111"),
		SubID:       7,
		LimitPrice:  decimal.RequireFromString("100"),
		Amount:      decimal.RequireFromString("0.5"),
		MaxFee:      decimal.RequireFromString("1"),
		RecipientID: 42,
		IsBid:       true,
	}
}

func TestTradeModuleData_Hash_Determinism(t *testing.T) {
	t1, err := baseTrade().Hash()
	require.NoError(t, err)
	t2, err := baseTrade().Hash()
	require.NoError(t, err)
	assert.Equal(t, t1, t2)
}

func TestTradeModuleData_Hash_SensitiveToIsBid(t *testing.T) {
	a := baseTrade()
	b := baseTrade()
	b.IsBid = false
	ha, err := a.Hash()
	require.NoError(t, err)
	hb, err := b.Hash()
	require.NoError(t, err)
	assert.NotEqual(t, ha, hb)
}

func TestTradeModuleData_Hash_SensitiveToAmount(t *testing.T) {
	a, b := baseTrade(), baseTrade()
	b.Amount = decimal.RequireFromString("1")
	ha, _ := a.Hash()
	hb, _ := b.Hash()
	assert.NotEqual(t, ha, hb)
}

func TestTradeModuleData_Hash_SensitiveToPrice(t *testing.T) {
	a, b := baseTrade(), baseTrade()
	b.LimitPrice = decimal.RequireFromString("101")
	ha, _ := a.Hash()
	hb, _ := b.Hash()
	assert.NotEqual(t, ha, hb)
}

func TestTradeModuleData_Hash_SensitiveToRecipient(t *testing.T) {
	a, b := baseTrade(), baseTrade()
	b.RecipientID = 99
	ha, _ := a.Hash()
	hb, _ := b.Hash()
	assert.NotEqual(t, ha, hb)
}

func TestTradeModuleData_Hash_RejectsNegativeMaxFee(t *testing.T) {
	t1 := baseTrade()
	t1.MaxFee = decimal.RequireFromString("-1")
	_, err := t1.Hash()
	assert.Error(t, err)
}

func TestTradeModuleData_Hash_RejectsTooPreciseAmount(t *testing.T) {
	t1 := baseTrade()
	t1.Amount = decimal.RequireFromString("0.0000000000000000005") // 19 dp
	_, err := t1.Hash()
	assert.Error(t, err)
}

func TestTradeModuleData_Hash_RejectsTooPrecisePrice(t *testing.T) {
	t1 := baseTrade()
	t1.LimitPrice = decimal.RequireFromString("0.00000000000000000099")
	_, err := t1.Hash()
	assert.Error(t, err)
}

func TestTransferModuleData_Hash_Happy(t *testing.T) {
	tm := auth.TransferModuleData{
		ToSubaccount: 99,
		Asset:        common.HexToAddress("0x1111111111111111111111111111111111111111"),
		SubID:        3,
		Amount:       decimal.RequireFromString("10"),
	}
	h, err := tm.Hash()
	require.NoError(t, err)
	// `h` is `[32]byte`; len is a compile-time constant. Verify the
	// hash is non-zero — the trivial "Hash() returns zero" failure mode.
	assert.NotEqual(t, [32]byte{}, h)
}

func TestTransferModuleData_Hash_AllowsNegativeAmount(t *testing.T) {
	// Transfer amount is signed (positions can transfer negative quantities).
	tm := auth.TransferModuleData{
		ToSubaccount: 1,
		Asset:        common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Amount:       decimal.RequireFromString("-5"),
	}
	_, err := tm.Hash()
	assert.NoError(t, err)
}

func TestTransferModuleData_Hash_RejectsTooPreciseAmount(t *testing.T) {
	tm := auth.TransferModuleData{
		ToSubaccount: 1,
		Asset:        common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Amount:       decimal.RequireFromString("0.0000000000000000005"),
	}
	_, err := tm.Hash()
	assert.Error(t, err)
}
