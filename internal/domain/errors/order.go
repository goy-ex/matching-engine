// Package errors collects domain-level sentinel values and structured error types
// for order and order-book validation failures.
package errors

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
