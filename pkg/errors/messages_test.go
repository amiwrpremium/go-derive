package errors_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestDescription_KnownCodeReturnsText(t *testing.T) {
	got := derrors.Description(derrors.CodeRateLimitExceeded)
	assert.NotEmpty(t, got)
	assert.Contains(t, got, "rate")
}

func TestDescription_UnknownCodeReturnsEmpty(t *testing.T) {
	assert.Equal(t, "", derrors.Description(99999))
}

func TestHasDescription_KnownCode(t *testing.T) {
	assert.True(t, derrors.HasDescription(derrors.CodeInvalidSignature))
}

func TestHasDescription_UnknownCode(t *testing.T) {
	assert.False(t, derrors.HasDescription(99999))
}

// TestDescription_AllCodesHaveText runs a description coverage check —
// every Code* constant declared in codes.go must have a non-empty entry
// in the message map. If a future PR adds a new code without a message,
// this test fails immediately.
func TestDescription_AllCodesHaveText(t *testing.T) {
	for name, code := range allCodes() {
		t.Run(name, func(t *testing.T) {
			got := derrors.Description(code)
			assert.NotEmpty(t, got, "code %d (%s) has no description", code, name)
		})
	}
}

func TestDescription_NoTrailingPunctuation(t *testing.T) {
	// Descriptions should read as fragments (no terminal period) so they
	// compose cleanly into log lines via fmt.Errorf("%v: %s", err, ctx).
	for name, code := range allCodes() {
		desc := derrors.Description(code)
		assert.False(t, strings.HasSuffix(desc, "."),
			"description for %s ends with a period: %q", name, desc)
	}
}

func TestDescription_LowercaseStart(t *testing.T) {
	for name, code := range allCodes() {
		desc := derrors.Description(code)
		if desc == "" {
			continue
		}
		first := desc[0]
		// "USDC", "X-LyraWallet" — uppercase prefixes that are proper nouns
		// or HTTP header names — are allowed. Otherwise the first character
		// should be lowercase to read like a fragment.
		if first >= 'A' && first <= 'Z' {
			// Spot-check: only allow our known proper-noun prefixes.
			allowed := []string{"USDC", "X-LyraWallet", "RFQ", "PMRM", "WebSocket", "ERC-20", "Swell"}
			ok := false
			for _, p := range allowed {
				if strings.HasPrefix(desc, p) {
					ok = true
					break
				}
			}
			assert.True(t, ok, "%s: description starts with uppercase but isn't a known proper noun: %q", name, desc)
		}
	}
}

func TestAPIError_Error_FillsInCanonicalWhenMessageEmpty(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeRateLimitExceeded}
	got := e.Error()
	assert.Contains(t, got, "rate")
	assert.Contains(t, got, "-32000")
}

func TestAPIError_Error_KeepsServerMessageWhenPresent(t *testing.T) {
	e := &derrors.APIError{
		Code:    derrors.CodeRateLimitExceeded,
		Message: "Custom server message",
	}
	got := e.Error()
	assert.Contains(t, got, "Custom server message")
	assert.NotContains(t, got, "rate limit window")
}

func TestAPIError_Error_UnknownCodeAndEmptyMessage(t *testing.T) {
	e := &APIError99999{}
	_ = e
	apiErr := &derrors.APIError{Code: 99999}
	got := apiErr.Error()
	// No canonical description, no server message — output still has the
	// code embedded.
	assert.Contains(t, got, "99999")
}

func TestAPIError_CanonicalMessage_KnownCode(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeMMPFrozen}
	assert.Contains(t, e.CanonicalMessage(), "market-maker protection")
}

func TestAPIError_CanonicalMessage_UnknownCode(t *testing.T) {
	e := &derrors.APIError{Code: 99999}
	assert.Equal(t, "", e.CanonicalMessage())
}

// dummy private type to silence unused-import diagnostics if test file
// reorganisation drops the apiError synonym; harmless.
type APIError99999 struct{}
