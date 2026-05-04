package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestOptionType_Valid_Call(t *testing.T) { assert.True(t, enums.OptionTypeCall.Valid()) }
func TestOptionType_Valid_Put(t *testing.T)  { assert.True(t, enums.OptionTypePut.Valid()) }

func TestOptionType_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "straddle", "CALL"} {
		assert.False(t, enums.OptionType(v).Valid(), "value %q", v)
	}
}
