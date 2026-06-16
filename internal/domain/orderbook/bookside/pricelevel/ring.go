package pricelevel

import (
	"github.com/akhmy/goy-ex-matching-engine/internal/domain"
	"github.com/akhmy/goy-ex-matching-engine/internal/domain/errors"
)

const minCap = 10
const expandCapMultiplier = 2
const shrinkCapDenominator = 2

type RingPriceLevel struct {
	buf  []*domain.Order
	head int
	tail int
	len  int
	cap  int
}

func NewRingPriceLevel() *RingPriceLevel {
	return &RingPriceLevel{
		buf:  make([]*domain.Order, minCap),
		head: 0,
		tail: 0,
		len:  0,
		cap:  minCap,
	}
}

func (pl *RingPriceLevel) Front() *domain.Order {
	if pl.len == 0 {
		return nil
	}
	return pl.buf[pl.head]
}

func (pl *RingPriceLevel) RemoveFront() {
	if pl.len == 0 {
		panic(errors.InvariantViolation(ErrRemoveFrontFromEmptyPriceLevel))
	}

	pl.buf[pl.head] = nil

	pl.head = (pl.head + 1) % pl.cap
	pl.len--

	if pl.len > 0 && pl.len <= pl.cap/4 && pl.cap > minCap {
		pl.shrink()
	}
}

func (pl *RingPriceLevel) PushBack(order *domain.Order) {
	if pl.len == pl.cap {
		pl.expand()
	}

	if pl.len > 0 {
		pl.tail = (pl.tail + 1) % pl.cap
	}

	pl.buf[pl.tail] = order
	pl.len++
}

func (pl *RingPriceLevel) Len() int {
	return pl.len
}

func (pl *RingPriceLevel) expand() {
	expandedCap := pl.cap * expandCapMultiplier
	expanded := make([]*domain.Order, expandedCap)

	for i := range pl.len {
		expanded[i] = pl.buf[(pl.head+i)%pl.cap]
	}

	pl.cap = expandedCap
	pl.buf = expanded
	pl.head = 0
	pl.tail = pl.len - 1
}

func (pl *RingPriceLevel) shrink() {
	shrunkCap := pl.cap / shrinkCapDenominator
	shrunk := make([]*domain.Order, shrunkCap)

	for i := range pl.len {
		shrunk[i] = pl.buf[(pl.head+i)%pl.cap]
	}

	pl.cap = shrunkCap
	pl.buf = shrunk
	pl.head = 0
	pl.tail = pl.len - 1
}
