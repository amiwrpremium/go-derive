package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestSentinels_DistinctValues(t *testing.T) {
	assert.NotEqual(t, derrors.ErrUnauthorized, derrors.ErrRateLimited)
	assert.NotEqual(t, derrors.ErrNotConnected, derrors.ErrAlreadyConnected)
	assert.NotEqual(t, derrors.ErrSubaccountRequired, derrors.ErrInvalidConfig)
	assert.NotEqual(t, derrors.ErrSubscriptionClosed, derrors.ErrNotConnected)
}

func TestSentinels_HaveMessages(t *testing.T) {
	for _, e := range []error{
		derrors.ErrNotConnected,
		derrors.ErrAlreadyConnected,
		derrors.ErrUnauthorized,
		derrors.ErrRateLimited,
		derrors.ErrSubscriptionClosed,
		derrors.ErrSubaccountRequired,
		derrors.ErrInvalidConfig,
	} {
		assert.NotEmpty(t, e.Error())
		assert.Contains(t, e.Error(), "derive")
	}
}

func TestExportedStdlibHelpers(t *testing.T) {
	assert.NotNil(t, derrors.Is)
	assert.NotNil(t, derrors.As)
	assert.NotNil(t, derrors.Unwrap)
	assert.NotNil(t, derrors.New)

	e := derrors.New("hello")
	assert.Equal(t, "hello", e.Error())
}
