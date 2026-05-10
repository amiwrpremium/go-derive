package ws

import "encoding/json"

// decodeJSON is the boilerplate JSON decoder every Derive channel
// uses — unmarshal the payload bytes into a typed T. The typed
// Subscribe* methods on [Client] all wire it up internally; the
// generic [Subscribe] is exposed for callers that need a custom
// decoder (e.g. a yet-undocumented channel that re-shapes its
// payload before delivery).
func decodeJSON[T any](raw json.RawMessage) (T, error) {
	var v T
	err := json.Unmarshal(raw, &v)
	return v, err
}
