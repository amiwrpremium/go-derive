package auth_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func baseAction() auth.ActionData {
	return auth.ActionData{
		SubaccountID: 1,
		Nonce:        42,
		Module:       common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Data:         [32]byte{0x42},
		Expiry:       1700000000,
		Owner:        common.HexToAddress("0x2222222222222222222222222222222222222222"),
		Signer:       common.HexToAddress("0x3333333333333333333333333333333333333333"),
	}
}

func TestActionData_Hash_Determinism(t *testing.T) {
	a := baseAction()
	assert.Equal(t, a.Hash(), a.Hash())
}

func TestActionData_Hash_SensitiveToSubaccountID(t *testing.T) {
	a, b := baseAction(), baseAction()
	b.SubaccountID = 2
	assert.NotEqual(t, a.Hash(), b.Hash())
}

func TestActionData_Hash_SensitiveToNonce(t *testing.T) {
	a, b := baseAction(), baseAction()
	b.Nonce = 43
	assert.NotEqual(t, a.Hash(), b.Hash())
}

func TestActionData_Hash_SensitiveToModule(t *testing.T) {
	a, b := baseAction(), baseAction()
	b.Module = common.HexToAddress("0x9999999999999999999999999999999999999999")
	assert.NotEqual(t, a.Hash(), b.Hash())
}

func TestActionData_Hash_SensitiveToData(t *testing.T) {
	a, b := baseAction(), baseAction()
	b.Data = [32]byte{0x99}
	assert.NotEqual(t, a.Hash(), b.Hash())
}

func TestActionData_Hash_SensitiveToExpiry(t *testing.T) {
	a, b := baseAction(), baseAction()
	b.Expiry++
	assert.NotEqual(t, a.Hash(), b.Hash())
}

func TestActionData_Hash_SensitiveToOwner(t *testing.T) {
	a, b := baseAction(), baseAction()
	b.Owner = common.HexToAddress("0x9999999999999999999999999999999999999999")
	assert.NotEqual(t, a.Hash(), b.Hash())
}

func TestActionData_Hash_SensitiveToSigner(t *testing.T) {
	a, b := baseAction(), baseAction()
	b.Signer = common.HexToAddress("0x9999999999999999999999999999999999999999")
	assert.NotEqual(t, a.Hash(), b.Hash())
}

func TestActionData_Hash_AllZerosIsValid(t *testing.T) {
	var a auth.ActionData
	h := a.Hash()
	assert.Len(t, h, 32)
}
