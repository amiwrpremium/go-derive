package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestEnvironment_Valid_Mainnet(t *testing.T) {
	assert.True(t, enums.EnvironmentMainnet.Valid())
}
func TestEnvironment_Valid_Testnet(t *testing.T) {
	assert.True(t, enums.EnvironmentTestnet.Valid())
}

func TestEnvironment_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "prod", "MAINNET", "staging"} {
		assert.False(t, enums.Environment(v).Valid(), "value %q", v)
	}
}

func TestEnvironment_Validate(t *testing.T) {
	assert.NoError(t, enums.EnvironmentMainnet.Validate())
	assert.NoError(t, enums.EnvironmentTestnet.Validate())

	err := enums.Environment("staging").Validate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, enums.ErrInvalidEnum)
	assert.Contains(t, err.Error(), "Environment")
	assert.Contains(t, err.Error(), "staging")
}
