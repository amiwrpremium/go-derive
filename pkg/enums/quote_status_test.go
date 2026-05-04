package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestQuoteStatus_Valid(t *testing.T) {
	for _, q := range []enums.QuoteStatus{
		enums.QuoteStatusOpen, enums.QuoteStatusFilled,
		enums.QuoteStatusCancelled, enums.QuoteStatusExpired,
	} {
		t.Run(string(q), func(t *testing.T) { assert.True(t, q.Valid()) })
	}
}

func TestQuoteStatus_Terminal(t *testing.T) {
	assert.False(t, enums.QuoteStatusOpen.Terminal())
	assert.True(t, enums.QuoteStatusFilled.Terminal())
	assert.True(t, enums.QuoteStatusCancelled.Terminal())
	assert.True(t, enums.QuoteStatusExpired.Terminal())
}

func TestQuoteStatus_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.QuoteStatus("").Valid())
	assert.False(t, enums.QuoteStatus("pending").Valid())
}
