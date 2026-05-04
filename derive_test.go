package goderive_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	goderive "github.com/amiwrpremium/go-derive"
)

func TestVersion_NonEmpty(t *testing.T) {
	assert.NotEmpty(t, goderive.Version)
}

func TestVersion_FormatLooksLikeSemver(t *testing.T) {
	// Either "X.Y.Z" or "X.Y.Z-suffix".
	assert.Regexp(t, `^\d+\.\d+\.\d+(-\w+)?$`, goderive.Version)
}

func TestUserAgent_StartsWithSDK(t *testing.T) {
	ua := goderive.UserAgent()
	assert.True(t, strings.HasPrefix(ua, "go-derive/"))
	assert.Contains(t, ua, goderive.Version)
}
