package netconf_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/internal/netconf"
)

func TestMainnet_Values(t *testing.T) {
	c := netconf.Mainnet()
	assert.Equal(t, netconf.NetworkMainnet, c.Network)
	assert.Equal(t, int64(957), c.ChainID)
	assert.True(t, strings.HasPrefix(c.HTTPURL, "https://"))
	assert.True(t, strings.HasPrefix(c.WSURL, "wss://"))
	assert.NotEmpty(t, c.Contracts.MatchingEngine)
	assert.NotEmpty(t, c.Contracts.TradeModule)
	assert.NotEmpty(t, c.Contracts.DepositModule)
	assert.NotEmpty(t, c.Contracts.WithdrawModule)
	assert.NotEmpty(t, c.Contracts.TransferModule)
}

func TestTestnet_Values(t *testing.T) {
	c := netconf.Testnet()
	assert.Equal(t, netconf.NetworkTestnet, c.Network)
	assert.Equal(t, int64(901), c.ChainID)
	assert.Contains(t, c.HTTPURL, "demo")
}

func TestMainnetTestnet_Distinct(t *testing.T) {
	m := netconf.Mainnet()
	te := netconf.Testnet()
	assert.NotEqual(t, m.ChainID, te.ChainID)
	assert.NotEqual(t, m.HTTPURL, te.HTTPURL)
	assert.NotEqual(t, m.WSURL, te.WSURL)
	assert.NotEqual(t, m.Network, te.Network)
}

// TestMainnetTestnet_SigningModulesDistinct guards against the
// copy-paste class of bug where a testnet module address ends up in
// the mainnet slot (or vice versa). The TradeModule and RFQModule are
// the two addresses the SDK actually hashes against today via
// signedOrderParams / signedQuoteParams; using the wrong network's
// address makes every signature reject server-side.
func TestMainnetTestnet_SigningModulesDistinct(t *testing.T) {
	m := netconf.Mainnet()
	te := netconf.Testnet()
	assert.NotEqual(t, m.Contracts.TradeModule, te.Contracts.TradeModule,
		"mainnet TradeModule must not equal testnet TradeModule")
	assert.NotEqual(t, m.Contracts.RFQModule, te.Contracts.RFQModule,
		"mainnet RFQModule must not equal testnet RFQModule")
}

func TestNetwork_String_Mainnet(t *testing.T) {
	assert.Equal(t, "mainnet", netconf.NetworkMainnet.String())
}

func TestNetwork_String_Testnet(t *testing.T) {
	assert.Equal(t, "testnet", netconf.NetworkTestnet.String())
}

func TestNetwork_String_Unknown(t *testing.T) {
	assert.Equal(t, "unknown(0)", netconf.NetworkUnknown.String())
}

func TestNetwork_String_OutOfRange(t *testing.T) {
	// Default arm of the switch.
	assert.Equal(t, "unknown(99)", netconf.Network(99).String())
	assert.Equal(t, "unknown(-1)", netconf.Network(-1).String())
}
