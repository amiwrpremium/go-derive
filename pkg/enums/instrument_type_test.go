package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestInstrumentType_Valid_Perp(t *testing.T) { assert.True(t, enums.InstrumentTypePerp.Valid()) }
func TestInstrumentType_Valid_Option(t *testing.T) {
	assert.True(t, enums.InstrumentTypeOption.Valid())
}
func TestInstrumentType_Valid_ERC20(t *testing.T) { assert.True(t, enums.InstrumentTypeERC20.Valid()) }

func TestInstrumentType_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "future", "spot", "PERP"} {
		assert.False(t, enums.InstrumentType(v).Valid(), "value %q", v)
	}
}
