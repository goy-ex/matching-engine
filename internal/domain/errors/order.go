// Package errors collects domain-level sentinel values and structured error types
// for order and order-book validation failures.
package errors

import (
	"errors"
	"fmt"
)

// InvalidSideError is returned when a Side value is neither SideBid nor SideAsk.
type InvalidSideError struct {
	Has string // the raw value that was rejected
}

func (e *InvalidSideError) Error() string {
	if e.Has == "" {
		e.Has = "[not provided]"
	}

	return "bad order side: " + e.Has
}

// ErrInvariantViolation is a sentinel used to wrap errors that represent a broken
// internal guarantee — a situation that should be impossible given correct usage.
var ErrInvariantViolation = errors.New("invariant violation")

// InvariantViolation wraps e together with ErrInvariantViolation so callers can
// detect the category with errors.Is while still seeing the specific cause.
func InvariantViolation(e error) error {
	return fmt.Errorf("%w: %w", ErrInvariantViolation, e)
}
