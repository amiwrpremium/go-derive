// Package auth.
package auth

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// HTTPHeaders builds the per-request authentication headers Derive expects
// on every REST call:
//
//	X-LyraWallet     — the owner address as 0x-prefixed hex
//	X-LyraTimestamp  — the current time as milliseconds since the Unix epoch
//	X-LyraSignature  — the EIP-191 signature over the timestamp string
//
// Despite the rename to "Derive", the header names retain their "Lyra"
// prefix server-side.
//
// If signer is nil, HTTPHeaders returns (nil, nil) — used by the public-only
// path of the HTTP transport. Errors from [Signer.SignAuthHeader] are
// propagated unmodified.
func HTTPHeaders(ctx context.Context, signer Signer, now time.Time) (http.Header, error) {
	if signer == nil {
		return nil, nil
	}
	sig, err := signer.SignAuthHeader(ctx, now)
	if err != nil {
		return nil, err
	}
	h := make(http.Header, 3)
	h.Set("X-LyraWallet", signer.Owner().Hex())
	h.Set("X-LyraTimestamp", strconv.FormatInt(now.UnixMilli(), 10))
	h.Set("X-LyraSignature", sig.Hex())
	return h, nil
}
