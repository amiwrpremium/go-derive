// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shape for the public operations /
// status endpoints — currently `public/get_live_incidents`.
package types

// Incident is one entry in `public/get_live_incidents.incidents`.
// Each entry is one ongoing platform incident (matched to an
// internal monitor or filed manually) with a severity and a
// human-readable message.
//
// The shape mirrors `IncidentResponseSchema` in Derive's v2.2
// OpenAPI spec.
//
// MonitorType ("manual" / "auto") and Severity ("low" / "medium" /
// "high") are bare strings for now; a later enum-tightening pass may
// type them.
type Incident struct {
	// CreationTimestampSec is when the incident was filed (Unix
	// seconds).
	CreationTimestampSec int64 `json:"creation_timestamp_sec"`
	// Label is the short incident label (e.g. "matching-engine").
	Label string `json:"label"`
	// Message is the longer human-readable message.
	Message string `json:"message"`
	// MonitorType is "manual" (filed by an operator) or "auto"
	// (raised by an internal monitor).
	MonitorType string `json:"monitor_type"`
	// Severity is "low", "medium", or "high".
	Severity string `json:"severity"`
}
