// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shape of `public/get_perp_impact_twap`.
package types

// PerpImpactTWAP is the response of `public/get_perp_impact_twap`. It
// reports time-weighted-average difference of mid, ask-impact, and
// bid-impact prices versus spot for one currency over the requested
// window.
//
// The shape mirrors `PublicGetPerpImpactTwapResultSchema` in
// `derivexyz/cockpit/orderbook-types`.
type PerpImpactTWAP struct {
	// Currency is the asset symbol (e.g. "BTC").
	Currency string `json:"currency"`
	// MidPriceDiffTWAP is the TWAP of (mid price − spot price) over
	// the window.
	MidPriceDiffTWAP Decimal `json:"mid_price_diff_twap"`
	// AskImpactDiffTWAP is the TWAP of (ask impact price − spot
	// price) over the window.
	AskImpactDiffTWAP Decimal `json:"ask_impact_diff_twap"`
	// BidImpactDiffTWAP is the TWAP of (bid impact price − spot
	// price) over the window.
	BidImpactDiffTWAP Decimal `json:"bid_impact_diff_twap"`
}
