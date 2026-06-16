package bookside

import "errors"

// ErrNoBestLevel is the panic argument when the best-price level is expected to
// exist but is missing from the internal map — indicates index/map desync.
var ErrNoBestLevel = errors.New("invariant violation: best level doesn't exist")

// ErrNoFront is the panic argument when the front element of a price-level list
// is nil despite the list being non-empty — indicates state corruption.
var ErrNoFront = errors.New("invariant violation: front element doesn't exist")

// ErrReduceBestOrderOnEmptyBookSide is the panic argument when ReduceBestOrder is
// called on a book side that has no price levels at all.
var ErrReduceBestOrderOnEmptyBookSide = errors.New("invariant violation: ReduceBestOrder called on empty book side")
