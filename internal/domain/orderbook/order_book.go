package orderbook

import (
	"github.com/goy-ex/matching-engine/internal/domain"
	"github.com/goy-ex/matching-engine/internal/pkg/sentinel"

	"github.com/google/uuid"
	"github.com/goy-ex/matching-engine/internal/domain/errors"
	"github.com/goy-ex/matching-engine/internal/domain/orderbook/bookside"
	"github.com/shopspring/decimal"
)

// BookSide is the interface that wraps a single half of the order book (bids or asks).
type BookSide interface {
	// GetBestOrder returns the highest-priority resting order on this side, or
	// false if the side is empty. For bids that is the highest price; for asks
	// the lowest price.
	GetBestOrder() (o *domain.Order, exists bool)

	// ReduceBestOrder subtracts amount from the best order's remaining quantity,
	// removing it (and its price level if empty) when fully filled.
	ReduceBestOrder(amount decimal.Decimal)

	// AddOrder inserts o into the correct price level, creating it if needed.
	AddOrder(o *domain.Order)
}

// OrderBook holds the bid and ask sides of the book for a single trading pair.
type OrderBook struct {
	Bids BookSide
	Asks BookSide
}

// New returns an empty OrderBook with sorted-slice book sides.
func New() *OrderBook {
	return &OrderBook{
		Bids: bookside.NewSortedSliceBookSide(true),
		Asks: bookside.NewSortedSliceBookSide(false),
	}
}

// Match runs the taker order against resting orders on the opposite side of the
// book. Each price level that crosses produces a Trade. Any unfilled remainder is
// rested on taker's own side. Returns all trades generated (may be empty).
func (ob *OrderBook) Match(taker *domain.Order) []*domain.Trade {
	var trades []*domain.Trade

	opposite := ob.oppositeSide(taker.Side)

	for taker.Remaining.GreaterThan(decimal.Zero) {
		maker, exists := opposite.GetBestOrder()
		if !exists {
			break
		}

		if !pricesCross(taker, maker) {
			break
		}

		matchAmount := decimal.Min(taker.Remaining, maker.Remaining)

		bidOrderID, askOrderID := resolveBidAsk(taker, maker)
		trade := domain.NewTrade(bidOrderID, askOrderID, maker.Price, matchAmount)
		trades = append(trades, trade)

		taker.Remaining = taker.Remaining.Sub(matchAmount)
		opposite.ReduceBestOrder(matchAmount)
	}

	if taker.Remaining.GreaterThan(decimal.Zero) {
		ob.oppositeSide(taker.Side.Opposite()).AddOrder(taker)
	}

	return trades
}

// oppositeSide returns the book side that holds resting orders against which the
// given side is matched: bids match against asks and vice versa.
func (ob *OrderBook) oppositeSide(side domain.Side) BookSide {
	switch side {
	case domain.SideBid:
		return ob.Asks
	case domain.SideAsk:
		return ob.Bids
	default:
		panic(sentinel.InvariantViolation(&errors.InvalidSideError{Has: string(side)}))
	}
}

func pricesCross(taker, maker *domain.Order) bool {
	switch taker.Side {
	case domain.SideBid:
		return taker.Price.GreaterThanOrEqual(maker.Price)
	case domain.SideAsk:
		return taker.Price.LessThanOrEqual(maker.Price)
	default:
		panic(sentinel.InvariantViolation(&errors.InvalidSideError{Has: string(taker.Side)}))
	}
}

//nolint:gocritic // ...
func resolveBidAsk(taker, maker *domain.Order) (uuid.UUID, uuid.UUID) {
	switch taker.Side {
	case domain.SideBid:
		return taker.ID, maker.ID
	case domain.SideAsk:
		return maker.ID, taker.ID
	default:
		panic(sentinel.InvariantViolation(&errors.InvalidSideError{Has: string(taker.Side)}))
	}
}
