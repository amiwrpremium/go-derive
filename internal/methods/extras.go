// Package methods is the shared implementation of every JSON-RPC method
// Derive exposes. Both pkg/rest.Client and pkg/ws.Client embed *API so that
// each method is defined exactly once, parameterised by the underlying
// transport.
//
// Public methods are unauthenticated; private methods require Signer to be
// non-nil. Private methods that mutate orders also use the Domain to sign
// the per-action EIP-712 hash.
package methods

// This file is now empty. Every JSON-RPC wrapper that used to live here
// was retyped against Derive's published v2.2 OpenAPI spec and moved to
// its domain-specific file:
//
//   - GetAccount, GetMargin, GetPublicMargin   → account.go
//   - GetMMPConfig                             → mmp.go
//   - GetFundingHistory, GetLiquidationHistory,
//     GetOptionSettlementHistory, GetSubaccountValueHistory,
//     GetERC20TransferHistory, GetInterestHistory,
//     ExpiredAndCancelledHistory,
//     GetPublicOptionSettlementHistory         → history.go
//   - GetNotifications, UpdateNotifications    → notifications.go
//   - Replace, OrderDebug, CancelByNonce,
//     SetCancelOnDisconnect                    → orders.go
//   - ChangeSubaccountLabel                    → subaccounts.go
//   - GetFundingRateHistory, GetSpotFeedHistory,
//     GetLatestSignedFeeds, GetPerpImpactTWAP  → feeds.go
//   - GetStatistics                            → markets.go
//   - GetTransaction                           → transactions.go
//
// The file is kept as a tombstone so accidental git-blame queries
// resolve. It will be removed in the cleanup commit.
