package domain

import (
	"strconv"

	"github.com/goy-ex/matching-engine/internal/domain/errors"
	"github.com/goy-ex/matching-engine/internal/pkg/sentinel"
)

// Side indicates which half of the order book an order belongs to.
type Side byte

const (
	// SideBid stands for buyer side: orders to purchase at or above a given price.
	SideBid Side = iota
	// SideAsk stands for seller side: orders to sell at or below a given price.
	SideAsk
)

// IsValid reports whether s is a recognized Side value.
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
		panic(sentinel.InvariantViolation(&errors.InvalidSideError{Has: string(s)}))
	}
}

func (s Side) String() string {
	switch s {
	case SideBid:
		return "bid"
	case SideAsk:
		return "ask"
	default:
		return strconv.Itoa(int(s))
	}
}
