// Package errors collects domain-level sentinel values and structured error types
// for order and order-book validation failures.
package errors

import "fmt"

// InvalidSideError is returned when a Side value is neither SideBid nor SideAsk.
type InvalidSideError struct {
	Has string // the raw value that was rejected
}

const defaultStringValue = "<non-provided>"

func (e *InvalidSideError) Error() string {
	if e.Has == "" {
		e.Has = defaultStringValue
	}

	return "bad order side: " + e.Has
}

type NegativeValueError struct {
	Field string
	Has   string
}

func (e *NegativeValueError) Error() string {
	if e.Has == "" {
		e.Has = defaultStringValue
	}

	if e.Field == "" {
		e.Field = defaultStringValue
	}

	return fmt.Sprintf("negative %s: %s", e.Field, e.Has)
}
