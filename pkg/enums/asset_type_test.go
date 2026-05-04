package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestAssetType_Valid(t *testing.T) {
	for _, a := range []enums.AssetType{
		enums.AssetTypeERC20, enums.AssetTypeOption, enums.AssetTypePerp,
	} {
		t.Run(string(a), func(t *testing.T) { assert.True(t, a.Valid()) })
	}
}

func TestAssetType_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.AssetType("").Valid())
	assert.False(t, enums.AssetType("future").Valid())
}
