package errors

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

// ErrNilOrderInPriceLevel is returned when a price-level list element holds a nil order pointer,
// indicating internal state corruption.
var ErrNilOrderInPriceLevel = errors.New("nil order in price level")

// UnexpectedTypeInPriceLevelError is returned when a list element in a price level contains
// a value that cannot be asserted to *domain.Order.
type UnexpectedTypeInPriceLevelError struct {
	Has any // the actual value found instead of *domain.Order
}

func (e *UnexpectedTypeInPriceLevelError) Error() string {
	return fmt.Sprintf("unexpected type in price level: %v", e.Has)
}

// ExistingEmptyPriceLevelError is returned when a price level with no orders is discovered
// while it should have been cleaned up after its last order was removed.
type ExistingEmptyPriceLevelError struct {
	Price decimal.Decimal
}

func (e *ExistingEmptyPriceLevelError) Error() string {
	if e.Price.IsZero() {
		return "empty price level exists for price: [not provided]"
	}

	return "empty price level exists for price: " + e.Price.String()
}
