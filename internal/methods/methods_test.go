package methods_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/methods"
	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

const testKey = "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"

// newAPI returns a *methods.API wired to a FakeTransport.
// signed=true attaches a LocalSigner and a default subaccount id.
func newAPI(t *testing.T, signed bool, sub int64) (*methods.API, *testutil.FakeTransport) {
	t.Helper()
	ft := testutil.NewFakeTransport()
	api := &methods.API{
		T:               ft,
		Domain:          netconf.Mainnet().EIP712Domain(),
		Nonces:          auth.NewNonceGen(),
		SignatureExpiry: 1_000,
	}
	api.SetTradeModule(common.HexToAddress(netconf.Mainnet().Contracts.TradeModule))
	if signed {
		s, err := auth.NewLocalSigner(testKey)
		require.NoError(t, err)
		api.Signer = s
		api.Subaccount = sub
	}
	return api, ft
}

func paramsAsMap(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()
	if len(raw) == 0 {
		return nil
	}
	var m map[string]any
	require.NoError(t, json.Unmarshal(raw, &m))
	return m
}
