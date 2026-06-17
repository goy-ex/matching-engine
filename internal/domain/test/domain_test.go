package domain_test

import (
	"math/rand/v2"
	"testing"

	"github.com/goy-ex/matching-engine/internal/domain/orderbook"
)

func TestPriceLevels(t *testing.T) {
	ob := orderbook.New()

	rand.Shuffle(len(smallFixture), func(i, j int) {
		smallFixture[i], smallFixture[j] = smallFixture[j], smallFixture[i]
	})

	for _, order := range smallFixture {
		ob.Match(order)
	}
}
