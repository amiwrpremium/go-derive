package enums_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestRFQInvalidReason_Valid(t *testing.T) {
	cases := []enums.RFQInvalidReason{
		enums.RFQInvalidReasonAccountUnderMaintenance,
		enums.RFQInvalidReasonWouldUnderMaintenance,
		enums.RFQInvalidReasonRiskReducingOnly,
		enums.RFQInvalidReasonReduceSize,
		enums.RFQInvalidReasonReduceOrCancel,
		enums.RFQInvalidReasonCancelLimitsOrUseIOC,
		enums.RFQInvalidReasonInsufficientBuyingPower,
	}
	for _, r := range cases {
		assert.True(t, r.Valid(), "expected %q valid", r)
		require.NoError(t, r.Validate())
	}
}

func TestRFQInvalidReason_Invalid(t *testing.T) {
	// Empty (the "no reason" wire value) is intentionally not
	// considered Valid by this method.
	assert.False(t, enums.RFQInvalidReason("").Valid())
	assert.False(t, enums.RFQInvalidReason("something else").Valid())
	err := enums.RFQInvalidReason("nope").Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, enums.ErrInvalidEnum))
}
