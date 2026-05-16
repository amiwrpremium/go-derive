package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestCancelReason_Valid_AllArms(t *testing.T) {
	cases := []enums.CancelReason{
		enums.CancelReasonNone,
		enums.CancelReasonUserRequest,
		enums.CancelReasonMMP,
		enums.CancelReasonInsufficientMargin,
		enums.CancelReasonSignedMaxFeeTooLow,
		enums.CancelReasonIOC,
		enums.CancelReasonCancelOnDisconnect,
		enums.CancelReasonSessionKey,
		enums.CancelReasonSubaccountWithdrawn,
		enums.CancelReasonCompliance,
		enums.CancelReasonTriggerFailed,
		enums.CancelReasonValidationFailed,
		enums.CancelReasonAlgoCompleted,
	}
	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			assert.True(t, c.Valid())
		})
	}
}

func TestCancelReason_Valid_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.CancelReason("nope").Valid())
	assert.False(t, enums.CancelReason("self_cross").Valid()) // hallucinated value
	assert.False(t, enums.CancelReason("expired").Valid())    // expired is OrderStatus, not CancelReason
}
