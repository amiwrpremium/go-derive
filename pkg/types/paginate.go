// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds [Paginate], the generic driver every "fetch every page"
// helper in the SDK delegates to.
package types

import "context"

// PaginateOptions tunes the [Paginate] primitive.
//
// All fields are optional. Pass a zero value to fetch every page using
// the server's default page size.
type PaginateOptions struct {
	// PageSize is the per-request page size threaded into each fetch's
	// [PageRequest]. Zero leaves PageSize zero and the server applies
	// its default (typically 100, max 1000).
	PageSize int
	// MaxItems caps the total number of records returned. When the
	// accumulator reaches MaxItems the result is truncated to that
	// length and pagination stops, even if more pages exist. Zero
	// means no cap.
	MaxItems int
	// MaxPages caps the number of pages fetched. Zero means no cap.
	MaxPages int
}

// Paginate drives a paginated SDK method to exhaustion (or until one of
// the caps in [PaginateOptions] is hit).
//
// The fetch closure receives a [PageRequest] (with Page and PageSize set
// per [PaginateOptions]) and must return one page of results plus the
// engine's [Page] envelope. The closure typically just calls the
// underlying SDK method, threading any query / filter args from the
// enclosing scope.
//
// Pagination stops when any of the following is true:
//
//   - The fetch returns an empty slice (defensive — some endpoints emit
//     NumPages=0 for "no records at all").
//   - We've just consumed the last page (Page >= NumPages).
//   - The accumulator reaches MaxItems (the result is truncated to
//     exactly MaxItems).
//   - We've consumed MaxPages pages.
//   - ctx.Err() is non-nil; that error is returned.
//
// Any error from the fetch is returned as-is; partial results
// accumulated up to that point are discarded.
//
// Example — fetch every order in a date window, capped at 1000:
//
//	all, err := types.Paginate(ctx, types.PaginateOptions{MaxItems: 1000},
//	    func(ctx context.Context, p types.PageRequest) ([]types.Order, types.Page, error) {
//	        return c.GetOrderHistory(ctx, q, p)
//	    })
func Paginate[T any](
	ctx context.Context,
	opts PaginateOptions,
	fetch func(ctx context.Context, page PageRequest) ([]T, Page, error),
) ([]T, error) {
	var all []T
	pageNum := 1
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		req := PageRequest{Page: pageNum, PageSize: opts.PageSize}
		results, info, err := fetch(ctx, req)
		if err != nil {
			return nil, err
		}
		if len(results) == 0 {
			break
		}
		all = append(all, results...)

		if opts.MaxItems > 0 && len(all) >= opts.MaxItems {
			all = all[:opts.MaxItems]
			break
		}
		if info.NumPages > 0 && pageNum >= info.NumPages {
			break
		}
		if opts.MaxPages > 0 && pageNum >= opts.MaxPages {
			break
		}
		pageNum++
	}
	return all, nil
}
