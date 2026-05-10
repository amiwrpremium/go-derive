// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the referrer-tree shapes returned by
// `public/get_descendant_tree` and `public/get_tree_roots`.
package types

import "encoding/json"

// DescendantTree is the response of `public/get_descendant_tree` —
// the tree of referees rooted at one wallet (or invite code).
//
// `Descendants` is preserved as raw JSON because the wire shape is
// recursive (each descendant is itself a tree); decode further at
// the call site if you need to walk it.
type DescendantTree struct {
	// Parent is the wallet (or invite code) at the root of the
	// returned subtree.
	Parent string `json:"parent"`
	// Descendants is the recursive subtree of referees. Inner shape
	// mirrors `DescendantTree` per the docs.
	Descendants json.RawMessage `json:"descendants"`
}

// TreeRoots is the response of `public/get_tree_roots` — the list
// of root wallets (top-of-tree referrers).
//
// `Roots` is preserved as raw JSON because the inner shape varies
// per program.
type TreeRoots struct {
	// Roots is the raw list of root entries.
	Roots json.RawMessage `json:"roots"`
}
