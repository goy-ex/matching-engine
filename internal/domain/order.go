package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/goy-ex/matching-engine/internal/domain/errors"
	"github.com/goy-ex/matching-engine/internal/pkg/sentinel"
	"github.com/shopspring/decimal"
)

// Order represents a single limit order placed by a user.
type Order struct {
	CreatedAt time.Time       // wall-clock time when the order was submitted
	Pair      string          // trading pair, e.g. "BTC/USDT"
	Price     decimal.Decimal // limit price the user is willing to trade at
	Amount    decimal.Decimal // original total quantity requested
	Remaining decimal.Decimal // quantity still unfilled; decremented as matches occur
	ID        uuid.UUID       // unique identifier of this order
	UserID    uuid.UUID       // identifier of the user who placed the order
	Side      Side            // whether this is a buy (Bid) or sell (Ask) order
}

// NewOrder validates the side and returns a new Order, or an error if side is invalid.
func NewOrder(
	id uuid.UUID,
	userID uuid.UUID,
	side Side,
	pair string,
	price decimal.Decimal,
	amount decimal.Decimal,
	remaining decimal.Decimal,
	createdAt time.Time,
) (*Order, error) {
	if !side.IsValid() {
		return nil, sentinel.BadRequest(&errors.InvalidSideError{Has: string(side)})
	}

	if price.LessThan(decimal.Zero) {
		return nil, sentinel.BadRequest(&errors.NegativeValueError{Field: "price", Has: price.String()})
	}

	if amount.LessThan(decimal.Zero) {
		return nil, sentinel.BadRequest(&errors.NegativeValueError{Field: "amount", Has: amount.String()})
	}

	if remaining.LessThan(decimal.Zero) {
		return nil, sentinel.BadRequest(&errors.NegativeValueError{Field: "remaining", Has: remaining.String()})
	}

	return &Order{
		ID:        id,
		UserID:    userID,
		Side:      side,
		Pair:      pair,
		Price:     price,
		Amount:    amount,
		Remaining: remaining,
		CreatedAt: createdAt,
	}, nil
}
