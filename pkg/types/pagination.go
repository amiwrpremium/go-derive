// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// All numeric fields use [Decimal], a thin wrapper around shopspring/decimal,
// so price/size/fee values never lose precision through float64 round-trips.
// On the wire, [Decimal] reads and writes JSON strings (Derive's preferred
// representation); a fallback path also accepts JSON numbers for resilience.
//
// Identifier types ([Address], [TxHash], [MillisTime]) carry the same
// round-trip guarantees: each one preserves the canonical wire format
// regardless of how Go marshals the surrounding struct.
//
// # Why named types
//
// Plain string and int64 fields would parse just fine, but named types let
// the SDK enforce invariants at construction time (NewAddress checksum
// check, NewDecimal precision check) and let callers tell at a glance which
// values are amounts vs prices vs subaccount ids.
package types

// Page wraps the server-side pagination shape used by every Derive list
// endpoint. The fields mirror the JSON response exactly — Derive returns
// just the totals and lets the caller track which page they asked for.
type Page struct {
	// NumPages is the total number of pages available for the query.
	NumPages int `json:"num_pages"`
	// Count is the total number of records across all pages.
	Count int `json:"count"`
}

// PageRequest is the common pagination input.
//
// Both Page and PageSize are 1-indexed; zero values are omitted on the wire
// so the server's defaults apply.
type PageRequest struct {
	// Page selects which page to fetch (1-indexed). Zero asks for the default.
	Page int `json:"page,omitempty"`
	// PageSize sets how many records per page. Zero asks for the default.
	PageSize int `json:"page_size,omitempty"`
}

// NewPageRequest constructs a [PageRequest] with both fields zero, which
// asks the server for its defaults.
func NewPageRequest() PageRequest { return PageRequest{} }

// WithPage returns a copy with the 1-indexed page set.
func (p PageRequest) WithPage(page int) PageRequest { p.Page = page; return p }

// WithPageSize returns a copy with the page size set.
func (p PageRequest) WithPageSize(size int) PageRequest { p.PageSize = size; return p }

// Validate enforces the schema: both fields must be non-negative.
// A zero in either slot is interpreted as "use the server default" by
// the json `omitempty` tag.
func (p PageRequest) Validate() error {
	if p.Page < 0 {
		return invalidParam("page", "must be non-negative")
	}
	if p.PageSize < 0 {
		return invalidParam("page_size", "must be non-negative")
	}
	return nil
}
