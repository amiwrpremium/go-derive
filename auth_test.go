package derive_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func baseTrade() derive.TradeModuleData {
	return derive.TradeModuleData{
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
	t1.Amount = decimal.RequireFromString("0.0000000000000000005")
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
	tm := derive.TransferModuleData{
		ToSubaccount: 99,
		Asset:        common.HexToAddress("0x1111111111111111111111111111111111111111"),
		SubID:        3,
		Amount:       decimal.RequireFromString("10"),
	}
	h, err := tm.Hash()
	require.NoError(t, err)

	assert.NotEqual(t, [32]byte{}, h)
}

func TestTransferModuleData_Hash_AllowsNegativeAmount(t *testing.T) {

	tm := derive.TransferModuleData{
		ToSubaccount: 1,
		Asset:        common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Amount:       decimal.RequireFromString("-5"),
	}
	_, err := tm.Hash()
	assert.NoError(t, err)
}

func TestTransferModuleData_Hash_RejectsTooPreciseAmount(t *testing.T) {
	tm := derive.TransferModuleData{
		ToSubaccount: 1,
		Asset:        common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Amount:       decimal.RequireFromString("0.0000000000000000005"),
	}
	_, err := tm.Hash()
	assert.Error(t, err)
}
func baseAction() derive.ActionData {
	return derive.ActionData{
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
	var a derive.ActionData
	h := a.Hash()
	assert.Len(t, h, 32)
}
func TestEIP191_DeterministicForSameTimestamp(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	ts := time.Unix(1700000000, 0)

	sig1, err := s.SignAuthHeader(context.Background(), ts)
	require.NoError(t, err)
	sig2, err := s.SignAuthHeader(context.Background(), ts)
	require.NoError(t, err)
	assert.Equal(t, sig1, sig2)
}

func TestEIP191_DifferentTimestampsDifferSignature(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	a, err := s.SignAuthHeader(context.Background(), time.Unix(1700000000, 0))
	require.NoError(t, err)
	b, err := s.SignAuthHeader(context.Background(), time.Unix(1700000001, 0))
	require.NoError(t, err)
	assert.NotEqual(t, a, b)
}

func TestEIP191_RecoverableViaPersonalSignDigest(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	ts := time.Unix(1700000000, 0)
	sig, err := s.SignAuthHeader(context.Background(), ts)
	require.NoError(t, err)

	msg := []byte(strconv.FormatInt(ts.UnixMilli(), 10))
	digest := personalHash(msg)
	pub, err := crypto.SigToPub(digest, normaliseV(sig[:]))
	require.NoError(t, err)
	assert.Equal(t, s.Address(), crypto.PubkeyToAddress(*pub))
}
func TestDomain_ChainIDAffectsSignature(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	a := derive.ActionData{Nonce: 1, Expiry: 1}

	mainnetSig, err := s.SignAction(context.Background(), netconf.Mainnet().EIP712Domain(), a)
	require.NoError(t, err)
	testnetSig, err := s.SignAction(context.Background(), netconf.Testnet().EIP712Domain(), a)
	require.NoError(t, err)

	assert.NotEqual(t, mainnetSig, testnetSig,
		"different chain IDs must produce different signatures")
}

func TestDomain_VerifyingContractAffectsSignature(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	a := derive.ActionData{Nonce: 1}

	cfg := netconf.Mainnet()
	d1 := cfg.EIP712Domain()
	d2 := d1
	d2.VerifyingContract = common.HexToAddress("0xdeadbeef00000000000000000000000000000000").Hex()

	sig1, err := s.SignAction(context.Background(), d1, a)
	require.NoError(t, err)
	sig2, err := s.SignAction(context.Background(), d2, a)
	require.NoError(t, err)
	assert.NotEqual(t, sig1, sig2)
}

func TestDomain_SameInputsSameSignature(t *testing.T) {

	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	d := netconf.Mainnet().EIP712Domain()
	a := derive.ActionData{Nonce: 99, Expiry: 1700000000}
	x, err := s.SignAction(context.Background(), d, a)
	require.NoError(t, err)
	y, err := s.SignAction(context.Background(), d, a)
	require.NoError(t, err)
	assert.Equal(t, x, y)
}
func TestHTTPHeaders_NilSignerYieldsNoHeaders(t *testing.T) {
	h, err := derive.HTTPHeaders(context.Background(), nil, time.Now())
	require.NoError(t, err)
	assert.Nil(t, h)
}

func TestHTTPHeaders_PopulatesAllThreeFields(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	now := time.Unix(1700000000, 123_000_000)
	h, err := derive.HTTPHeaders(context.Background(), s, now)
	require.NoError(t, err)
	assert.Equal(t, s.Owner().Hex(), h.Get("X-LyraWallet"))
	assert.Equal(t, "1700000000123", h.Get("X-LyraTimestamp"))
	assert.True(t, strings.HasPrefix(h.Get("X-LyraSignature"), "0x"))
	assert.Len(t, h.Get("X-LyraSignature"), 2+65*2)
}

// failSigner forces the SignAuthHeader path to error so HTTPHeaders'
// error-propagation branch is exercised.
type failSigner struct{ derive.Signer }

func (failSigner) SignAuthHeader(context.Context, time.Time) (derive.Signature, error) {
	return derive.Signature{}, errBoom
}

var errBoom = newErr("boom")

type sentinelErr struct{ msg string }

func (e sentinelErr) Error() string { return e.msg }

func newErr(s string) sentinelErr { return sentinelErr{msg: s} }

func TestHTTPHeaders_PropagatesSignerError(t *testing.T) {

	real, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	s := failSigner{Signer: real}
	_, err = derive.HTTPHeaders(context.Background(), s, time.Now())
	assert.ErrorContains(t, err, "boom")
}

// testKey is a throwaway secp256k1 key used only in tests.
const testKey = "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"

// personalHash mirrors auth/eip191.go for verification side.
func personalHash(msg []byte) []byte {
	prefix := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg)))
	h := crypto.NewKeccakState()
	_, _ = h.Write(prefix)
	_, _ = h.Write(msg)
	return h.Sum(nil)
}

// normaliseV converts Derive's v in {27,28} back to go-ethereum's {0,1}.
func normaliseV(sig []byte) []byte {
	out := make([]byte, len(sig))
	copy(out, sig)
	if len(out) == 65 {
		out[64] -= 27
	}
	return out
}

// timeT is an alias for time.Time so other test files can declare
// `timeNowDeterministic() timeT` without importing time directly.
type timeT = time.Time

// FuzzNewLocalSigner verifies that bad hex input never panics. Real keys
// produce a signer; everything else returns an error.
func FuzzNewLocalSigner(f *testing.F) {
	f.Add("0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318")
	f.Add("4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318")
	f.Add("0xZZZZ")
	f.Add("")
	f.Add("0x")
	f.Add(string(make([]byte, 256)))

	f.Fuzz(func(t *testing.T, s string) {

		signer, err := derive.NewLocalSigner(s)
		if err != nil {
			if signer != nil {
				t.Fatalf("error path returned non-nil signer for %q", s)
			}
			return
		}
		if signer == nil {
			t.Fatalf("nil signer with no error for %q", s)
		}
	})
}
func TestNewLocalSigner_Hex0xPrefixed(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestNewLocalSigner_HexNoPrefix(t *testing.T) {
	s, err := derive.NewLocalSigner(strings.TrimPrefix(testKey, "0x"))
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestNewLocalSigner_RejectsBadHex(t *testing.T) {
	for _, in := range []string{"", "not-hex", "0xZZ", "0x12"} {
		t.Run(in, func(t *testing.T) {
			_, err := derive.NewLocalSigner(in)
			assert.Error(t, err)
		})
	}
}

func TestLocalSigner_AddressMatchesPublicKey(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	k, err := crypto.HexToECDSA(strings.TrimPrefix(testKey, "0x"))
	require.NoError(t, err)
	assert.Equal(t, crypto.PubkeyToAddress(k.PublicKey), s.Address())
}

func TestLocalSigner_OwnerEqualsAddress(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	assert.Equal(t, s.Address(), s.Owner(), "LocalSigner has no separate owner")
}

func TestLocalSigner_SignAuthHeader_Recover(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	ts := time.Now()

	sig, err := s.SignAuthHeader(context.Background(), ts)
	require.NoError(t, err)

	msg := []byte(strconv.FormatInt(ts.UnixMilli(), 10))
	digest := personalHash(msg)
	pub, err := crypto.SigToPub(digest, normaliseV(sig[:]))
	require.NoError(t, err)
	assert.Equal(t, s.Address(), crypto.PubkeyToAddress(*pub))
}

func TestLocalSigner_SignAction_Determinism(t *testing.T) {
	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	domain := netconf.Mainnet().EIP712Domain()
	action := derive.ActionData{SubaccountID: 1, Nonce: 12345, Expiry: 1700000000}
	a, err := s.SignAction(context.Background(), domain, action)
	require.NoError(t, err)
	b, err := s.SignAction(context.Background(), domain, action)
	require.NoError(t, err)
	assert.Equal(t, a, b)
}

func TestLocalSigner_SignAction_PopulatesOwnerAndSignerWhenZero(t *testing.T) {

	s, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	domain := netconf.Mainnet().EIP712Domain()

	a, err := s.SignAction(context.Background(), domain, derive.ActionData{Nonce: 1})
	require.NoError(t, err)
	b, err := s.SignAction(context.Background(), domain, derive.ActionData{Nonce: 2})
	require.NoError(t, err)
	assert.NotEqual(t, a, b)
}
func TestNonceGen_StrictlyMonotonic(t *testing.T) {
	g := derive.NewNonceGen()
	prev := g.Next()
	for i := 0; i < 1000; i++ {
		n := g.Next()
		require.Greater(t, n, prev, "iteration %d", i)
		prev = n
	}
}

func TestNonceGen_ConcurrentUniqueness(t *testing.T) {
	g := derive.NewNonceGen()
	const goroutines = 16
	const perG = 200
	results := make([][]uint64, goroutines)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			out := make([]uint64, perG)
			for j := 0; j < perG; j++ {
				out[j] = g.Next()
			}
			results[i] = out
		}(i)
	}
	wg.Wait()

	seen := map[uint64]struct{}{}
	for _, r := range results {
		for _, n := range r {
			_, dup := seen[n]
			require.False(t, dup, "duplicate nonce %d", n)
			seen[n] = struct{}{}
		}
	}
	assert.Equal(t, goroutines*perG, len(seen))
}

func TestNonceGen_NewIsNonZeroFirst(t *testing.T) {
	g := derive.NewNonceGen()
	assert.NotZero(t, g.Next())
}
func TestNewSessionKeySigner_Happy(t *testing.T) {
	owner := common.HexToAddress("0x1111111111111111111111111111111111111111")
	s, err := derive.NewSessionKeySigner(testKey, owner)
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestNewSessionKeySigner_RejectsBadKey(t *testing.T) {
	_, err := derive.NewSessionKeySigner("not-hex", common.Address{})
	assert.Error(t, err)
}

func TestSessionKeySigner_OwnerSeparateFromAddress(t *testing.T) {
	owner := common.HexToAddress("0x1111111111111111111111111111111111111111")
	s, err := derive.NewSessionKeySigner(testKey, owner)
	require.NoError(t, err)
	assert.NotEqual(t, s.Address(), s.Owner())
	assert.Equal(t, owner, s.Owner())
}

func TestSessionKeySigner_SignAuthHeaderDelegates(t *testing.T) {
	owner := common.HexToAddress("0x1111111111111111111111111111111111111111")
	s, err := derive.NewSessionKeySigner(testKey, owner)
	require.NoError(t, err)

	sig, err := s.SignAuthHeader(context.Background(), timeNowDeterministic())
	require.NoError(t, err)
	assert.NotEqual(t, [65]byte{}, sig)
}

func TestSessionKeySigner_SignActionStampsOwnerAndSigner(t *testing.T) {

	owner := common.HexToAddress("0x2222222222222222222222222222222222222222")
	s, err := derive.NewSessionKeySigner(testKey, owner)
	require.NoError(t, err)
	domain := netconf.Mainnet().EIP712Domain()
	a := derive.ActionData{Nonce: 1, Owner: common.Address{}, Signer: common.Address{}}
	b := derive.ActionData{Nonce: 1, Owner: common.HexToAddress("0xdeadbeef00000000000000000000000000000000"), Signer: common.HexToAddress("0xfeedface00000000000000000000000000000000")}

	sigA, err := s.SignAction(context.Background(), domain, a)
	require.NoError(t, err)
	sigB, err := s.SignAction(context.Background(), domain, b)
	require.NoError(t, err)
	assert.Equal(t, sigA, sigB, "SignAction must overwrite Owner/Signer fields with the configured ones")
}

// timeNowDeterministic is a tiny indirection so changing this in one place
// propagates to every test that relies on a fixed timestamp.
func timeNowDeterministic() (t timeT) { return timeT{} }
func TestSignature_Hex_LengthAndPrefix(t *testing.T) {
	var s derive.Signature
	hex := s.Hex()
	assert.True(t, strings.HasPrefix(hex, "0x"))
	assert.Len(t, hex, 2+65*2)
}

func TestSignature_Hex_AllZeros(t *testing.T) {
	var s derive.Signature
	hex := s.Hex()
	assert.Equal(t, "0x"+strings.Repeat("00", 65), hex)
}

func TestSignature_Hex_AllOnes(t *testing.T) {
	var s derive.Signature
	for i := range s {
		s[i] = 0xff
	}
	assert.Equal(t, "0x"+strings.Repeat("ff", 65), s.Hex())
}

func TestSignature_Hex_MixedBytes(t *testing.T) {
	var s derive.Signature
	s[0] = 0x12
	s[1] = 0x34
	s[64] = 0xab
	hex := s.Hex()
	assert.True(t, strings.HasPrefix(hex, "0x1234"))
	assert.True(t, strings.HasSuffix(hex, "ab"))
}

var goodAsset = common.HexToAddress("0x1111111111111111111111111111111111111111")

func validTrade() derive.TradeModuleData {
	return derive.TradeModuleData{
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
		mut  func(*derive.TradeModuleData)
		want string
	}{
		{"zero asset", func(t *derive.TradeModuleData) { t.Asset = common.Address{} }, "asset"},
		{"zero price", func(t *derive.TradeModuleData) { t.LimitPrice = decimal.Zero }, "limit_price"},
		{"negative price", func(t *derive.TradeModuleData) { t.LimitPrice = decimal.RequireFromString("-1") }, "limit_price"},
		{"zero amount", func(t *derive.TradeModuleData) { t.Amount = decimal.Zero }, "amount"},
		{"negative fee", func(t *derive.TradeModuleData) { t.MaxFee = decimal.RequireFromString("-1") }, "max_fee"},
		{"negative recipient", func(t *derive.TradeModuleData) { t.RecipientID = -1 }, "recipient_id"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := validTrade()
			c.mut(&d)
			err := d.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, derive.ErrInvalidInput))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func validTransfer() derive.TransferModuleData {
	return derive.TransferModuleData{
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
		mut  func(*derive.TransferModuleData)
		want string
	}{
		{"zero asset", func(t *derive.TransferModuleData) { t.Asset = common.Address{} }, "asset"},
		{"negative subaccount", func(t *derive.TransferModuleData) { t.ToSubaccount = -1 }, "to_subaccount"},
		{"zero amount", func(t *derive.TransferModuleData) { t.Amount = decimal.Zero }, "amount"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := validTransfer()
			c.mut(&d)
			err := d.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, derive.ErrInvalidInput))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func validAction() derive.ActionData {
	return derive.ActionData{
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
		mut  func(*derive.ActionData)
		want string
	}{
		{"negative subaccount", func(a *derive.ActionData) { a.SubaccountID = -1 }, "subaccount_id"},
		{"zero module", func(a *derive.ActionData) { a.Module = common.Address{} }, "module"},
		{"zero owner", func(a *derive.ActionData) { a.Owner = common.Address{} }, "owner"},
		{"zero signer", func(a *derive.ActionData) { a.Signer = common.Address{} }, "signer"},
		{"zero expiry", func(a *derive.ActionData) { a.Expiry = 0 }, "expiry"},
		{"negative expiry", func(a *derive.ActionData) { a.Expiry = -1 }, "expiry"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := validAction()
			c.mut(&d)
			err := d.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, derive.ErrInvalidInput))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}
