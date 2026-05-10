package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestAlgoType_Valid(t *testing.T) {
	assert.True(t, enums.AlgoTypeTWAP.Valid())
}

func TestAlgoType_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.AlgoType("").Valid())
	assert.False(t, enums.AlgoType("vwap").Valid())
}
