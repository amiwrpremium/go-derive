package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestLiquidityRole_Valid_Maker(t *testing.T) { assert.True(t, enums.LiquidityRoleMaker.Valid()) }
func TestLiquidityRole_Valid_Taker(t *testing.T) { assert.True(t, enums.LiquidityRoleTaker.Valid()) }

func TestLiquidityRole_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "MAKER", "neither"} {
		assert.False(t, enums.LiquidityRole(v).Valid(), "value %q", v)
	}
}
