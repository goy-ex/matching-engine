package domain_test

import (
	"time"

	"github.com/akhmy/goy-ex-matching-engine/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func mustOrder(side domain.Side, price, amount int64) *domain.Order {
	order, err := domain.NewOrder(
		uuid.New(),
		uuid.New(),
		side,
		"BTC/USDT",
		decimal.NewFromInt(price),
		decimal.NewFromInt(amount),
		decimal.NewFromInt(amount),
		time.Now(),
	)
	if err != nil {
		panic(err)
	}

	return order
}

var smallFixture = []*domain.Order{
	// bids (buy), highest to lowest
	mustOrder(domain.SideBid, 99_900, 1),
	mustOrder(domain.SideBid, 99_800, 2),
	mustOrder(domain.SideBid, 99_700, 3),
	mustOrder(domain.SideBid, 99_600, 1),
	mustOrder(domain.SideBid, 99_500, 5),
	mustOrder(domain.SideBid, 99_400, 2),
	mustOrder(domain.SideBid, 99_300, 4),
	mustOrder(domain.SideBid, 99_200, 1),
	mustOrder(domain.SideBid, 99_100, 3),
	mustOrder(domain.SideBid, 99_000, 6),
	// asks (sell), lowest to highest
	mustOrder(domain.SideAsk, 100_100, 2),
	mustOrder(domain.SideAsk, 100_200, 1),
	mustOrder(domain.SideAsk, 100_300, 3),
	mustOrder(domain.SideAsk, 100_400, 2),
	mustOrder(domain.SideAsk, 100_500, 4),
	mustOrder(domain.SideAsk, 100_600, 1),
	mustOrder(domain.SideAsk, 100_700, 5),
	mustOrder(domain.SideAsk, 100_800, 2),
	mustOrder(domain.SideAsk, 100_900, 3),
	mustOrder(domain.SideAsk, 101_000, 1),
}
