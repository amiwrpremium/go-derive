//go:build integration

package integration_test

import (
	"sort"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
)

// TestCross_GetInstruments verifies REST and WS return the same instrument
// set for the same query — they speak the same JSON-RPC method name on
// the wire, so they should agree.
func TestCross_GetInstruments(t *testing.T) {
	env := loadEnv(t)
	rest := newRESTClient(t, env)
	ws := newWSClient(t, env)

	ctx, cancel := withTimeout(t)
	defer cancel()

	restIns, err := rest.GetInstruments(ctx, "BTC", derive.InstrumentTypePerp)
	require.NoError(t, err)
	wsIns, err := ws.GetInstruments(ctx, "BTC", derive.InstrumentTypePerp)
	require.NoError(t, err)

	require.Equal(t, len(restIns), len(wsIns), "REST and WS should return the same number of instruments")

	restNames := make([]string, len(restIns))
	for i, in := range restIns {
		restNames[i] = in.Name
	}
	wsNames := make([]string, len(wsIns))
	for i, in := range wsIns {
		wsNames[i] = in.Name
	}
	sort.Strings(restNames)
	sort.Strings(wsNames)
	assert.Equal(t, restNames, wsNames)
}

// TestCross_GetTicker verifies REST and WS return tickers for the same
// instrument with marks within 1% of each other (tolerance accounts for
// the small lag between the two reads).
func TestCross_GetTicker(t *testing.T) {
	env := loadEnv(t)
	rest := newRESTClient(t, env)
	ws := newWSClient(t, env)

	ctx, cancel := withTimeout(t)
	defer cancel()

	restTk, err := rest.GetTicker(ctx, env.instrument)
	require.NoError(t, err)
	wsTk, err := ws.GetTicker(ctx, env.instrument)
	require.NoError(t, err)

	assert.Equal(t, restTk.InstrumentName, wsTk.InstrumentName)

	restMark := restTk.MarkPrice.Inner()
	wsMark := wsTk.MarkPrice.Inner()
	if restMark.IsZero() || wsMark.IsZero() {
		return
	}
	tolerance := restMark.Mul(decimal.RequireFromString("0.01"))
	delta := restMark.Sub(wsMark).Abs()
	assert.True(t, delta.Cmp(tolerance) <= 0,
		"REST mark %s and WS mark %s differ by more than 1%%", restMark, wsMark)
}

// TestCross_FacadeWiring exercises the [github.com/amiwrpremium/go-derive.Client] facade against the live
// network so a regression in option-threading shows up here.
func TestCross_FacadeWiring(t *testing.T) {
	env := loadEnv(t)
	c := newDeriveClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	insts, err := c.REST.GetInstruments(ctx, "BTC", derive.InstrumentTypePerp)
	require.NoError(t, err)
	assert.NotEmpty(t, insts)
}
