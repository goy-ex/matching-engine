// Package pricelevel defines the PriceLevel abstraction and its doubly-linked-list
// implementation. A price level is a FIFO queue of orders sharing the same limit
// price; orders are matched in arrival order (time priority).
package pricelevel

import (
	"container/list"

	"github.com/akhmy/goy-ex-matching-engine/internal/domain"
	"github.com/akhmy/goy-ex-matching-engine/internal/domain/errors"
)

type ListPriceLevel struct {
	orders *list.List
	len    int
}

// Front returns the first order in the level and panics if the stored value is not *domain.Order.
func (pl *ListPriceLevel) Front() *domain.Order {
	front := pl.orders.Front()
	if front == nil {
		return nil
	}

	order, ok := front.Value.(*domain.Order)
	if !ok {
		panic(errors.InvariantViolation(&InvalidPriceLevelElement{Has: front.Value}))
	}

	return order
}

func (pl *ListPriceLevel) RemoveFront() {
	front := pl.orders.Front()
	if front == nil {
		panic(errors.InvariantViolation(ErrRemoveFrontFromEmptyPriceLevel))
	}

	pl.orders.Remove(front)
	pl.len--
}

func (pl *ListPriceLevel) PushBack(order *domain.Order) {
	pl.orders.PushBack(order)
	pl.len++
}

// Len returns the number of orders currently in the level.
func (pl *ListPriceLevel) Len() int {
	return pl.len
}
