package derive_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive"
)

func TestMainnet_Values(t *testing.T) {
	c := derive.Mainnet()
	assert.Equal(t, derive.NetworkMainnet, c.Network)
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
	c := derive.Testnet()
	assert.Equal(t, derive.NetworkTestnet, c.Network)
	assert.Equal(t, int64(901), c.ChainID)
	assert.Contains(t, c.HTTPURL, "demo")
}

func TestMainnetTestnet_Distinct(t *testing.T) {
	m := derive.Mainnet()
	te := derive.Testnet()
	assert.NotEqual(t, m.ChainID, te.ChainID)
	assert.NotEqual(t, m.HTTPURL, te.HTTPURL)
	assert.NotEqual(t, m.WSURL, te.WSURL)
	assert.NotEqual(t, m.Network, te.Network)
}

func TestNetwork_String_Mainnet(t *testing.T) {
	assert.Equal(t, "mainnet", derive.NetworkMainnet.String())
}

func TestNetwork_String_Testnet(t *testing.T) {
	assert.Equal(t, "testnet", derive.NetworkTestnet.String())
}

func TestNetwork_String_Unknown(t *testing.T) {
	assert.Equal(t, "unknown(0)", derive.NetworkUnknown.String())
}

func TestNetwork_String_OutOfRange(t *testing.T) {
	// Default arm of the switch.
	assert.Equal(t, "unknown(99)", derive.Network(99).String())
	assert.Equal(t, "unknown(-1)", derive.Network(-1).String())
}

func TestEIP712Domain_BindsToConfig(t *testing.T) {
	c := derive.Mainnet()
	d := c.EIP712Domain()
	assert.Equal(t, "Matching", d.Name)
	assert.Equal(t, "1", d.Version)
	assert.Equal(t, c.ChainID, d.ChainID)
	assert.Equal(t, c.Contracts.MatchingEngine, d.VerifyingContract)
}

func TestEIP712Domain_DiffersBetweenNetworks(t *testing.T) {
	m := derive.Mainnet().EIP712Domain()
	te := derive.Testnet().EIP712Domain()
	assert.NotEqual(t, m.ChainID, te.ChainID)
	assert.NotEqual(t, m.VerifyingContract, te.VerifyingContract)
	// Name and Version are pinned, same on both.
	assert.Equal(t, m.Name, te.Name)
	assert.Equal(t, m.Version, te.Version)
}

func TestEIP712Domain_FromCustomConfig(t *testing.T) {
	c := derive.NetworkConfig{ChainID: 42, Contracts: derive.Contracts{MatchingEngine: "0xdeadbeef"}}
	d := c.EIP712Domain()
	assert.Equal(t, int64(42), d.ChainID)
	assert.Equal(t, "0xdeadbeef", d.VerifyingContract)
}
