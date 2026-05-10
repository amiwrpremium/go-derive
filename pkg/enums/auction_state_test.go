package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestAuctionState_Valid(t *testing.T) {
	assert.True(t, enums.AuctionStateOngoing.Valid())
	assert.True(t, enums.AuctionStateEnded.Valid())
}

func TestAuctionState_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.AuctionState("").Valid())
	assert.False(t, enums.AuctionState("paused").Valid())
}
