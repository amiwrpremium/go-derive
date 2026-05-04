package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestAuctionType_Valid(t *testing.T) {
	assert.True(t, enums.AuctionTypeSolvent.Valid())
	assert.True(t, enums.AuctionTypeInsolvent.Valid())
}

func TestAuctionType_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.AuctionType("").Valid())
	assert.False(t, enums.AuctionType("partial").Valid())
}
