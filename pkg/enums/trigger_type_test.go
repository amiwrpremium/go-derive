package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestTriggerType_Valid(t *testing.T) {
	assert.True(t, enums.TriggerTypeStopLoss.Valid())
	assert.True(t, enums.TriggerTypeTakeProfit.Valid())
}

func TestTriggerType_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.TriggerType("").Valid())
	assert.False(t, enums.TriggerType("trailing").Valid())
}

func TestTriggerPriceType_Valid(t *testing.T) {
	assert.True(t, enums.TriggerPriceTypeMark.Valid())
	assert.True(t, enums.TriggerPriceTypeIndex.Valid())
}

func TestTriggerPriceType_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.TriggerPriceType("").Valid())
	assert.False(t, enums.TriggerPriceType("last").Valid())
}
