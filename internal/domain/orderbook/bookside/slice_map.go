package bookside

import (
	"sort"

	"github.com/akhmy/goy-ex-matching-engine/internal/domain"
	"github.com/akhmy/goy-ex-matching-engine/internal/domain/orderbook/bookside/pricelevel"
	"github.com/akhmy/goy-ex-matching-engine/internal/pkg/sentinel"
	"github.com/shopspring/decimal"
)

// PriceLevel is a FIFO queue of limit orders at a single price point.
type priceLevel interface {
	// Front returns the oldest (highest time-priority) order, or nil if empty.
	Front() *domain.Order

	RemoveFront()

	PushBack(o *domain.Order)

	Len() int
}

// SortedSliceBookSide is a book-side implementation that keeps price levels in a
// sorted slice (desc=true for bids, desc=false for asks) and in a map for O(1)
// lookup by price. The best price is always at index 0 of the slice.
type SortedSliceBookSide struct {
	levels map[string]priceLevel
	prices []decimal.Decimal
	desc   bool
}

// NewSortedSliceBookSide returns an empty book side. Pass desc=true for the bid
// side (highest price first) and desc=false for the ask side (lowest price first).
func NewSortedSliceBookSide(desc bool) *SortedSliceBookSide {
	return &SortedSliceBookSide{
		levels: make(map[string]priceLevel),
		prices: make([]decimal.Decimal, 0),
		desc:   desc,
	}
}

// GetBestOrder returns the highest-priority resting order: highest-priced bid or
// lowest-priced ask. Returns false if the book side is empty.
func (bs *SortedSliceBookSide) GetBestOrder() (*domain.Order, bool) {
	if len(bs.prices) == 0 {
		return nil, false
	}

	level, exists := bs.levels[bs.prices[0].String()]
	if !exists {
		panic(sentinel.InvariantViolation(ErrNoBestLevel))
	}

	return level.Front(), true
}

// ReduceBestOrder subtracts amount from the best order's remaining quantity.
// If the order is fully filled it is removed; if its price level becomes empty
// the level is pruned. Panics with ErrReduceBestOrderOnEmptyBookSide if called
// on an empty book side.
func (bs *SortedSliceBookSide) ReduceBestOrder(amount decimal.Decimal) {
	if len(bs.prices) == 0 {
		panic(sentinel.InvariantViolation(ErrReduceBestOrderOnEmptyBookSide))
	}

	level, exists := bs.levels[bs.prices[0].String()]
	if !exists {
		panic(sentinel.InvariantViolation(ErrReduceBestOrderOnEmptyBookSide))
	}

	order := level.Front()
	if order == nil {
		panic(sentinel.InvariantViolation(ErrNoFront))
	}

	order.Remaining = order.Remaining.Sub(amount)
	if !order.Remaining.IsZero() {
		return
	}

	level.RemoveFront()

	if level.Len() > 0 {
		return
	}

	delete(bs.levels, bs.prices[0].String())
	bs.removePrice(order.Price)
}

// AddOrder inserts order into its price level, creating a new level if one does
// not yet exist for order.Price. The sorted prices slice is updated to keep the
// best price at index 0.
func (bs *SortedSliceBookSide) AddOrder(order *domain.Order) {
	if level, exists := bs.levels[order.Price.String()]; exists {
		level.PushBack(order)

		return
	}

	priceIndex := sort.Search(len(bs.prices), func(i int) bool {
		return bs.isBetter(order.Price, bs.prices[i])
	})

	bs.prices = append(bs.prices, decimal.Decimal{})
	copy(bs.prices[priceIndex+1:], bs.prices[priceIndex:])
	bs.prices[priceIndex] = order.Price

	newLevel := pricelevel.NewRingPriceLevel()
	newLevel.PushBack(order)
	bs.levels[order.Price.String()] = newLevel
}

func (bs *SortedSliceBookSide) isBetter(a, b decimal.Decimal) bool {
	if bs.desc {
		return a.GreaterThan(b)
	}

	return a.LessThan(b)
}

func (bs *SortedSliceBookSide) isBetterOrEqual(a, b decimal.Decimal) bool {
	if bs.desc {
		return a.GreaterThanOrEqual(b)
	}

	return a.LessThanOrEqual(b)
}

func (bs *SortedSliceBookSide) removePrice(price decimal.Decimal) {
	idx := sort.Search(len(bs.prices), func(i int) bool {
		return bs.isBetterOrEqual(price, bs.prices[i])
	})

	if idx < len(bs.prices) && bs.prices[idx].Equal(price) {
		copy(bs.prices[idx:], bs.prices[idx+1:])
		bs.prices = bs.prices[:len(bs.prices)-1]
	}
}
