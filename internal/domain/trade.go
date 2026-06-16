package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Trade records a matched execution between a bid order and an ask order.
type Trade struct {
	ExecutedAt time.Time       // wall-clock time when the match occurred
	Price      decimal.Decimal // maker's price at which the trade was executed
	Amount     decimal.Decimal // quantity exchanged
	ID         uuid.UUID       // unique identifier of this trade
	BidOrderID uuid.UUID       // ID of the buy-side order involved in the trade
	AskOrderID uuid.UUID       // ID of the sell-side order involved in the trade
}

// NewTrade creates a Trade stamped with the current time and a new random ID.
func NewTrade(
	bidOrderID uuid.UUID,
	askOrderID uuid.UUID,
	price decimal.Decimal,
	amount decimal.Decimal,
) *Trade {
	return &Trade{
		ID:         uuid.New(),
		BidOrderID: bidOrderID,
		AskOrderID: askOrderID,
		Price:      price,
		Amount:     amount,
		ExecutedAt: time.Now(),
	}
}
