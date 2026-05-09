package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestIncident_Decode(t *testing.T) {
	raw := []byte(`{
		"creation_timestamp_sec": 1700000000,
		"label": "matching-engine",
		"message": "Elevated matching latency on perp markets.",
		"monitor_type": "auto",
		"severity": "medium"
	}`)
	var i types.Incident
	require.NoError(t, json.Unmarshal(raw, &i))
	assert.Equal(t, int64(1700000000), i.CreationTimestampSec)
	assert.Equal(t, "matching-engine", i.Label)
	assert.Equal(t, "auto", i.MonitorType)
	assert.Equal(t, "medium", i.Severity)
}
