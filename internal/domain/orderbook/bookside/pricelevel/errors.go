package pricelevel

import (
	"errors"
	"fmt"
)

type InvalidPriceLevelElement struct {
	Has any
}

func (e *InvalidPriceLevelElement) Error() string {
	return fmt.Sprintf("invalid price element: %v", e.Has)
}

var ErrRemoveFrontFromEmptyPriceLevel = errors.New("RemoveFront called on empty price level")
