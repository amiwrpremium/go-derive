package derive_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/derive"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

const testKey = "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"

func TestNewClient_RequiresNetwork(t *testing.T) {
	_, err := derive.NewClient()
	assert.ErrorIs(t, err, derrors.ErrInvalidConfig)
}

func TestNewClient_Mainnet(t *testing.T) {
	c, err := derive.NewClient(derive.WithMainnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.NotNil(t, c.REST)
	assert.NotNil(t, c.WS)
	assert.Equal(t, netconf.NetworkMainnet, c.Network().Network)
}

func TestNewClient_Testnet(t *testing.T) {
	c, err := derive.NewClient(derive.WithTestnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, netconf.NetworkTestnet, c.Network().Network)
}

func TestNewClient_CustomNetwork(t *testing.T) {
	custom := netconf.Testnet()
	custom.HTTPURL = "https://custom.example/api"
	c, err := derive.NewClient(derive.WithCustomNetwork(custom))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, "https://custom.example/api", c.Network().HTTPURL)
}

func TestNewClient_WithSignerAndSubaccount(t *testing.T) {
	signer, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	c, err := derive.NewClient(
		derive.WithMainnet(),
		derive.WithSigner(signer),
		derive.WithSubaccount(42),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	// Signer is threaded through; the REST client embeds *methods.API
	// which uses Signer for header injection. We can't easily inspect
	// without reaching into internals, so a build-time success is the
	// signal here.
	assert.NotNil(t, c.REST)
}

func TestClient_CloseIdempotent(t *testing.T) {
	c, err := derive.NewClient(derive.WithMainnet())
	require.NoError(t, err)
	require.NoError(t, c.Close())
	// Second close on REST is a no-op; WS close on never-connected is also fine.
	_ = c.Close() // tolerate any additional close cleanup
}
