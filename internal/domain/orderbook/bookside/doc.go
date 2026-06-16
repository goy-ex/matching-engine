// Package bookside implements one half of an order book (either bids or asks).
// The central type, SortedSliceBookSide, keeps price levels sorted so that the
// best price is always at index 0 and can be accessed in O(1) time.
package bookside
