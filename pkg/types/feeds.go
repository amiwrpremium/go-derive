// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the public oracle-feed shapes returned by Derive's
// signed-feed endpoints: `public/get_funding_rate_history`,
// `public/get_spot_feed_history`, and `public/get_latest_signed_feeds`.
package types

// FundingRateHistoryItem is one entry in
// `public/get_funding_rate_history`. The endpoint reports the engine's
// hourly perpetual funding rate prints over the requested window.
//
// The shape mirrors `FundingRateSchema` in Derive's v2.2 OpenAPI
// spec.
type FundingRateHistoryItem struct {
	// Timestamp is the funding tick (millisecond Unix epoch).
	Timestamp MillisTime `json:"timestamp"`
	// FundingRate is the hourly funding rate at the tick.
	FundingRate Decimal `json:"funding_rate"`
}

// SpotFeedHistoryItem is one entry in `public/get_spot_feed_history`.
// The endpoint reports historical oracle spot prices for one currency
// at the requested period.
//
// The shape mirrors `SpotFeedHistoryResponseSchema` in Derive's v2.2
// OpenAPI spec.
type SpotFeedHistoryItem struct {
	// Timestamp is the sample point (millisecond Unix epoch).
	Timestamp MillisTime `json:"timestamp"`
	// TimestampBucket is the sample's bucket-start time (millisecond
	// Unix epoch). For aggregated history the engine stamps each
	// sample with the bucket it belongs to so callers can resample
	// without rewriting timestamps.
	TimestampBucket MillisTime `json:"timestamp_bucket"`
	// Price is the oracle price for the bucket.
	Price Decimal `json:"price"`
}

// SignedFeeds is the response of `public/get_latest_signed_feeds`. It
// returns the latest oracle-signed snapshot for every published
// currency / expiry, plus the signers and signatures the on-chain
// feed contracts will check.
//
// All five maps are keyed by currency at the outer level; the inner
// keys vary (feed type for `PerpData`, expiry for `FwdData`,
// `RateData`, and `VolData`). `SpotData` is a single layer because
// spot prices are per-currency, not per-expiry.
//
// The shape mirrors `PublicGetLatestSignedFeedsResultSchema`.
type SignedFeeds struct {
	// SpotData maps currency → latest spot-feed snapshot.
	SpotData map[string]SpotFeedData `json:"spot_data"`
	// PerpData maps currency → feed type ("P" / "A" / "B") → latest
	// perp-feed snapshot.
	PerpData map[string]map[string]PerpFeedData `json:"perp_data"`
	// FwdData maps currency → expiry (string-formatted Unix seconds)
	// → latest forward-feed snapshot.
	FwdData map[string]map[string]ForwardFeedData `json:"fwd_data"`
	// RateData maps currency → expiry → latest rate-feed snapshot.
	RateData map[string]map[string]RateFeedData `json:"rate_data"`
	// VolData maps currency → expiry → latest vol-feed snapshot.
	VolData map[string]map[string]VolFeedData `json:"vol_data"`
}

// SpotFeedData is one signed spot-price snapshot.
type SpotFeedData struct {
	// Currency is the asset symbol (e.g. "BTC").
	Currency string `json:"currency"`
	// Price is the oracle price.
	Price Decimal `json:"price"`
	// Confidence is the confidence interval the oracle attaches.
	Confidence Decimal `json:"confidence"`
	// Timestamp is when the snapshot was produced (millisecond Unix
	// epoch).
	Timestamp MillisTime `json:"timestamp"`
	// Deadline is the latest time the snapshot can be submitted
	// on-chain (millisecond Unix epoch).
	Deadline MillisTime `json:"deadline"`
	// FeedSourceType identifies the source: "S" for the staking
	// network, "O" for the optimistic feed. Optional on the wire;
	// default is "S".
	FeedSourceType string `json:"feed_source_type,omitempty"`
	// Signatures carries the signers and their signatures over the
	// snapshot.
	Signatures OracleSignatureData `json:"signatures"`
}

// PerpFeedData is one signed perp-feed snapshot.
type PerpFeedData struct {
	// Currency is the asset symbol.
	Currency string `json:"currency"`
	// Type is the perp feed flavour: "P" (mid), "A" (ask impact),
	// or "B" (bid impact).
	Type string `json:"type"`
	// SpotDiffValue is the difference between the perp mark and the
	// spot at the snapshot time.
	SpotDiffValue Decimal `json:"spot_diff_value"`
	// Confidence is the confidence interval the oracle attaches.
	Confidence Decimal `json:"confidence"`
	// Timestamp is when the snapshot was produced.
	Timestamp MillisTime `json:"timestamp"`
	// Deadline is the latest time the snapshot can be submitted
	// on-chain.
	Deadline MillisTime `json:"deadline"`
	// Signatures carries the signers and their signatures.
	Signatures OracleSignatureData `json:"signatures"`
}

// ForwardFeedData is one signed forward-curve snapshot.
type ForwardFeedData struct {
	// Currency is the asset symbol.
	Currency string `json:"currency"`
	// Expiry is the option/forward expiry (Unix seconds).
	Expiry int64 `json:"expiry"`
	// FwdDiff is the forward-spot difference at the snapshot time.
	FwdDiff Decimal `json:"fwd_diff"`
	// SpotAggregateLatest is the latest spot aggregate.
	SpotAggregateLatest Decimal `json:"spot_aggregate_latest"`
	// SpotAggregateStart is the spot aggregate at the window start.
	SpotAggregateStart Decimal `json:"spot_aggregate_start"`
	// Confidence is the confidence interval the oracle attaches.
	Confidence Decimal `json:"confidence"`
	// Timestamp is when the snapshot was produced.
	Timestamp MillisTime `json:"timestamp"`
	// Deadline is the latest time the snapshot can be submitted
	// on-chain.
	Deadline MillisTime `json:"deadline"`
	// Signatures carries the signers and their signatures.
	Signatures OracleSignatureData `json:"signatures"`
}

// RateFeedData is one signed risk-free-rate snapshot.
type RateFeedData struct {
	// Currency is the asset symbol.
	Currency string `json:"currency"`
	// Expiry is the option/forward expiry (Unix seconds).
	Expiry int64 `json:"expiry"`
	// Rate is the rate value.
	Rate Decimal `json:"rate"`
	// Confidence is the confidence interval the oracle attaches.
	Confidence Decimal `json:"confidence"`
	// Timestamp is when the snapshot was produced.
	Timestamp MillisTime `json:"timestamp"`
	// Deadline is the latest time the snapshot can be submitted
	// on-chain.
	Deadline MillisTime `json:"deadline"`
	// Signatures carries the signers and their signatures.
	Signatures OracleSignatureData `json:"signatures"`
}

// VolFeedData is one signed volatility-surface snapshot. The
// `VolData` field carries the SVI parameter set used to fit the
// surface for the given currency / expiry.
type VolFeedData struct {
	// Currency is the asset symbol.
	Currency string `json:"currency"`
	// Expiry is the option expiry (Unix seconds).
	Expiry int64 `json:"expiry"`
	// VolData is the SVI parameter set for the surface.
	VolData VolSVIParam `json:"vol_data"`
	// Confidence is the confidence interval the oracle attaches.
	Confidence Decimal `json:"confidence"`
	// Timestamp is when the snapshot was produced.
	Timestamp MillisTime `json:"timestamp"`
	// Deadline is the latest time the snapshot can be submitted
	// on-chain.
	Deadline MillisTime `json:"deadline"`
	// Signatures carries the signers and their signatures.
	Signatures OracleSignatureData `json:"signatures"`
}

// VolSVIParam is the SVI parameter set carried by a vol-feed
// snapshot. The parameters fit Gatheral's stochastic-volatility-
// inspired (SVI) form to the surface at the snapshot's expiry.
//
// Field names mirror the OAS verbatim — Derive's wire format uses
// underscored capitals (`SVI_a`, `SVI_b`, …).
type VolSVIParam struct {
	// SVIA is the level term.
	SVIA Decimal `json:"SVI_a"`
	// SVIB is the slope term.
	SVIB Decimal `json:"SVI_b"`
	// SVIFwd is the forward price reference.
	SVIFwd Decimal `json:"SVI_fwd"`
	// SVIM is the moneyness shift term.
	SVIM Decimal `json:"SVI_m"`
	// SVIRefTau is the reference time-to-expiry.
	SVIRefTau Decimal `json:"SVI_refTau"`
	// SVIRho is the correlation term.
	SVIRho Decimal `json:"SVI_rho"`
	// SVISigma is the at-the-money volatility term.
	SVISigma Decimal `json:"SVI_sigma"`
}

// OracleSignatureData carries the signers and their signatures over
// a feed snapshot. On-chain feed contracts verify a configurable
// quorum of these signatures before accepting the snapshot.
type OracleSignatureData struct {
	// Signers is the list of signer addresses.
	Signers []string `json:"signers,omitempty"`
	// Signatures is the matching list of signatures (1:1 with
	// Signers).
	Signatures []string `json:"signatures,omitempty"`
}
