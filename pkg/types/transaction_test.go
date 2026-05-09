package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestDepositTx_Decode(t *testing.T) {
	payload := `{
		"tx_hash": "0x1111111111111111111111111111111111111111111111111111111111111111",
		"asset": "USDC",
		"amount": "1000",
		"subaccount_id": 1,
		"status": "completed",
		"timestamp": 1700000000000
	}`
	var d types.DepositTx
	require.NoError(t, json.Unmarshal([]byte(payload), &d))
	assert.Equal(t, "USDC", d.Asset)
	assert.Equal(t, "completed", d.Status)
}

func TestWithdrawTx_Decode(t *testing.T) {
	payload := `{
		"tx_hash": "0x2222222222222222222222222222222222222222222222222222222222222222",
		"asset": "USDC",
		"amount": "500",
		"subaccount_id": 1,
		"status": "pending",
		"timestamp": 1700000000000
	}`
	var w types.WithdrawTx
	require.NoError(t, json.Unmarshal([]byte(payload), &w))
	assert.Equal(t, "pending", w.Status)
}

func TestDepositTx_RoundTrip(t *testing.T) {
	in := types.DepositTx{Asset: "USDC", Amount: types.MustDecimal("1"), Status: "completed"}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.DepositTx
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.Asset, out.Asset)
	assert.Equal(t, in.Status, out.Status)
}

func TestTransaction_Decode(t *testing.T) {
	raw := []byte(`{
		"data":"{\"foo\":\"bar\"}",
		"error_log":null,
		"status":"settled",
		"transaction_hash":"0xabc"
	}`)
	var got types.Transaction
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, enums.TxStatusSettled, got.Status)
	assert.Equal(t, "0xabc", got.TransactionHash)
	assert.Equal(t, "", got.ErrorLog)
}

func TestTransaction_FailedTx(t *testing.T) {
	// A failed transaction reports error_log; transaction_hash may be null.
	raw := []byte(`{
		"data":"...",
		"error_log":"reverted: insufficient balance",
		"status":"reverted",
		"transaction_hash":null
	}`)
	var got types.Transaction
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, enums.TxStatusReverted, got.Status)
	assert.Equal(t, "reverted: insufficient balance", got.ErrorLog)
	assert.Equal(t, "", got.TransactionHash)
}
