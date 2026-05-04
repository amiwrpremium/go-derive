package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestMarginType_Valid(t *testing.T) {
	for _, m := range []enums.MarginType{
		enums.MarginTypeSM, enums.MarginTypePM, enums.MarginTypePM2,
	} {
		t.Run(string(m), func(t *testing.T) { assert.True(t, m.Valid()) })
	}
}

func TestMarginType_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.MarginType("").Valid())
	assert.False(t, enums.MarginType("pm3").Valid())
	assert.False(t, enums.MarginType("sm").Valid()) // case-sensitive — must be uppercase
}
