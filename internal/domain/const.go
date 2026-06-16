package domain

import "github.com/akhmy/goy-ex-matching-engine/internal/domain/errors"

// Side indicates which half of the order book an order belongs to.
type Side byte

const (
	SideBid Side = iota // buyer side: orders to purchase at or above a given price
	SideAsk             // seller side: orders to sell at or below a given price
)

// IsValid reports whether s is a recognised Side value.
func (s Side) IsValid() bool {
	switch s {
	case SideBid, SideAsk:
		return true
	default:
		return false
	}
}

// Opposite returns the other side of the book (Bid→Ask, Ask→Bid).
// Panics with InvalidSideError if s is not a valid Side.
func (s Side) Opposite() Side {
	switch s {
	case SideBid:
		return SideAsk
	case SideAsk:
		return SideBid
	default:
		panic(&errors.InvalidSideError{Has: string(s)})
	}
}
