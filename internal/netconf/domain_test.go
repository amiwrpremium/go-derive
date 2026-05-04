package netconf_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/internal/netconf"
)

func TestEIP712Domain_BindsToConfig(t *testing.T) {
	c := netconf.Mainnet()
	d := c.EIP712Domain()
	assert.Equal(t, "Matching", d.Name)
	assert.Equal(t, "1", d.Version)
	assert.Equal(t, c.ChainID, d.ChainID)
	assert.Equal(t, c.Contracts.MatchingEngine, d.VerifyingContract)
}

func TestEIP712Domain_DiffersBetweenNetworks(t *testing.T) {
	m := netconf.Mainnet().EIP712Domain()
	te := netconf.Testnet().EIP712Domain()
	assert.NotEqual(t, m.ChainID, te.ChainID)
	assert.NotEqual(t, m.VerifyingContract, te.VerifyingContract)
	// Name and Version are pinned, same on both.
	assert.Equal(t, m.Name, te.Name)
	assert.Equal(t, m.Version, te.Version)
}

func TestEIP712Domain_FromCustomConfig(t *testing.T) {
	c := netconf.Config{ChainID: 42, Contracts: netconf.Contracts{MatchingEngine: "0xdeadbeef"}}
	d := c.EIP712Domain()
	assert.Equal(t, int64(42), d.ChainID)
	assert.Equal(t, "0xdeadbeef", d.VerifyingContract)
}
