package types

import (
	"encoding/json"
	"strconv"
	"time"
)

// MillisTime is a time.Time that round-trips as integer milliseconds since
// the Unix epoch — Derive's preferred timestamp format on every JSON-RPC
// payload.
//
// The zero value is the zero time.Time; use [time.Time.IsZero] on the
// underlying [MillisTime.Time] to detect it.
type MillisTime struct {
	// T is the underlying time. Use [MillisTime.Time] in callers; this
	// field is exported only so that struct literals are convenient.
	T time.Time
}

// NewMillisTime wraps a [time.Time] as a [MillisTime].
func NewMillisTime(t time.Time) MillisTime { return MillisTime{T: t} }

// Time returns the underlying [time.Time].
func (m MillisTime) Time() time.Time { return m.T }

// Millis returns the time as milliseconds since the Unix epoch.
func (m MillisTime) Millis() int64 { return m.T.UnixMilli() }

// MarshalJSON encodes the time as an integer count of milliseconds since
// the Unix epoch.
func (m MillisTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.T.UnixMilli())
}

// UnmarshalJSON decodes either a JSON number or a JSON string of integer
// milliseconds. Empty strings and JSON null leave the receiver as the zero
// value.
func (m *MillisTime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			return nil
		}
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		m.T = time.UnixMilli(n)
		return nil
	}
	var n int64
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	m.T = time.UnixMilli(n)
	return nil
}
