// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the stDRV (staked DRV) snapshot shape returned by
// `public/get_stdrv_snapshots`.
package types

// StDRVSnapshot is one staked-DRV balance snapshot at a single
// point in time, returned in the `snapshots` slice of
// `public/get_stdrv_snapshots`.
type StDRVSnapshot struct {
	// Amount is the staked-DRV balance at the snapshot.
	Amount Decimal `json:"amount"`
	// TimestampSec is the snapshot's Unix-seconds timestamp.
	TimestampSec int64 `json:"timestamp_sec"`
}

// StDRVSnapshots is the response of `public/get_stdrv_snapshots` —
// one wallet's staked-DRV balance over a time window.
type StDRVSnapshots struct {
	// Wallet is the queried wallet address.
	Wallet string `json:"wallet"`
	// Snapshots is the slice of balance snapshots, oldest first.
	Snapshots []StDRVSnapshot `json:"snapshots"`
}
