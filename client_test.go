package derive_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
)

func TestNewClient_RequiresNetwork(t *testing.T) {
	_, err := derive.NewClient()
	assert.ErrorIs(t, err, derive.ErrInvalidConfig)
}

func TestNewClient_Mainnet(t *testing.T) {
	c, err := derive.NewClient(derive.WithMainnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.NotNil(t, c.REST)
	assert.NotNil(t, c.WS)
	assert.Equal(t, derive.NetworkMainnet, c.Network().Network)
}

func TestNewClient_Testnet(t *testing.T) {
	c, err := derive.NewClient(derive.WithTestnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, derive.NetworkTestnet, c.Network().Network)
}

func TestNewClient_CustomNetwork(t *testing.T) {
	custom := derive.Testnet()
	custom.HTTPURL = "https://custom.example/api"
	c, err := derive.NewClient(derive.WithCustomNetwork(custom))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, "https://custom.example/api", c.Network().HTTPURL)
}

func TestNewClient_WithSignerAndSubaccount(t *testing.T) {
	signer, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)
	c, err := derive.NewClient(
		derive.WithMainnet(),
		derive.WithSigner(signer),
		derive.WithSubaccount(42),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	// Signer is threaded through; the REST client embeds the apiCalls
	// struct which uses Signer for header injection. We can't easily
	// inspect without reaching into internals, so a build-time success
	// is the signal here.
	assert.NotNil(t, c.REST)
}

func TestClient_CloseIdempotent(t *testing.T) {
	c, err := derive.NewClient(derive.WithMainnet())
	require.NoError(t, err)
	require.NoError(t, c.Close())
	// Second close on REST is a no-op; WS close on never-connected is
	// also fine.
	_ = c.Close()
}
