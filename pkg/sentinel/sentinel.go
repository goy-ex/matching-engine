// Package sentinel defines application-wide error sentinels and their wrapping
// helpers. Sentinels are used with errors.Is to classify errors at handler
// boundaries without coupling callers to concrete error types.
package sentinel

import (
	"errors"
	"fmt"
)

// ErrInvariantViolation indicates that an internal guarantee has been broken —
// a state that correct code should never reach. Typically results in a 500.
var ErrInvariantViolation = errors.New("invariant violation")

// InvariantViolation wraps e with ErrInvariantViolation so the category is
// detectable via errors.Is while the specific cause remains visible.
func InvariantViolation(e error) error {
	return fmt.Errorf("%w: %w", ErrInvariantViolation, e)
}

// ErrBadRequest indicates that the caller supplied invalid input.
// Typically results in a 400.
var ErrBadRequest = errors.New("bad request")

// BadRequest wraps e with ErrBadRequest so the category is detectable via
// errors.Is while the specific cause remains visible.
func BadRequest(e error) error {
	return fmt.Errorf("%w: %w", ErrBadRequest, e)
}
