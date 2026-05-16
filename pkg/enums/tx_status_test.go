package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestTxStatus_Valid(t *testing.T) {
	for _, s := range []enums.TxStatus{
		enums.TxStatusRequested, enums.TxStatusPending, enums.TxStatusSettled,
		enums.TxStatusReverted, enums.TxStatusIgnored, enums.TxStatusTimedOut,
	} {
		t.Run(string(s), func(t *testing.T) { assert.True(t, s.Valid()) })
	}
}

func TestTxStatus_Terminal(t *testing.T) {
	assert.False(t, enums.TxStatusRequested.Terminal())
	assert.False(t, enums.TxStatusPending.Terminal())
	assert.True(t, enums.TxStatusSettled.Terminal())
	assert.True(t, enums.TxStatusReverted.Terminal())
	assert.True(t, enums.TxStatusIgnored.Terminal())
	assert.True(t, enums.TxStatusTimedOut.Terminal())
}

func TestTxStatus_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.TxStatus("").Valid())
	assert.False(t, enums.TxStatus("done").Valid())
}
