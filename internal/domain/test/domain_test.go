package domain_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/goy-ex/matching-engine/internal/domain"
	domainerrors "github.com/goy-ex/matching-engine/internal/domain/errors"
	"github.com/shopspring/decimal"
)

func TestNewOrder(t *testing.T) {
	type orderInput struct {
		id        uuid.UUID
		userID    uuid.UUID
		side      domain.Side
		pair      string
		price     decimal.Decimal
		amount    decimal.Decimal
		remaining decimal.Decimal
		createdAt time.Time
	}

	validInput := orderInput{
		id:        uuid.New(),
		userID:    uuid.New(),
		side:      domain.Side(0),
		pair:      "BTC/USDT",
		price:     decimal.NewFromInt(1000),
		amount:    decimal.NewFromInt(5),
		remaining: decimal.NewFromInt(5),
		createdAt: time.Now(),
	}

	type testCase struct {
		order  orderInput
		assert func(error) error
	}

	testCases := map[string]testCase{
		"valid": {
			order: validInput,
			assert: func(err error) error {
				if err == nil {
					return nil
				}

				return fmt.Errorf("want nil, got: %v", err)
			},
		},
		"invalid_side": {
			order: func() orderInput {
				input := validInput
				input.side = domain.Side(10)

				return input
			}(),
			assert: func(err error) error {
				_, ok := errors.AsType[*domainerrors.InvalidSideError](err)
				if !ok {
					return fmt.Errorf("want InvalidSideError, got: %v", err)
				}

				return nil
			},
		},
		"negative_price": {
			order: func() orderInput {
				input := validInput
				input.price = decimal.NewFromInt(-100)

				return input
			}(),
			assert: func(err error) error {
				_, ok := errors.AsType[*domainerrors.NegativeValueError](err)
				if !ok {
					return fmt.Errorf("want NegativeValueError, got: %v", err)
				}

				return nil
			},
		},
		"negative_amount": {
			order: func() orderInput {
				input := validInput
				input.amount = decimal.NewFromInt(-100)

				return input
			}(),
			assert: func(err error) error {
				_, ok := errors.AsType[*domainerrors.NegativeValueError](err)
				if !ok {
					return fmt.Errorf("want NegativeValueError, got: %v", err)
				}

				return nil
			},
		},
		"negative_remaining": {
			order: func() orderInput {
				input := validInput
				input.remaining = decimal.NewFromInt(-100)

				return input
			}(),
			assert: func(err error) error {
				_, ok := errors.AsType[*domainerrors.NegativeValueError](err)
				if !ok {
					return fmt.Errorf("want NegativeValueError, got: %v", err)
				}

				return nil
			},
		},
	}

	for tcName, tc := range testCases {
		t.Run(tcName, func(t *testing.T) {
			t.Parallel()

			_, err := domain.NewOrder(
				tc.order.id,
				tc.order.userID,
				tc.order.side,
				tc.order.pair,
				tc.order.price,
				tc.order.amount,
				tc.order.remaining,
				tc.order.createdAt,
			)

			err = tc.assert(err)
			if err != nil {
				t.Fatalf("assert error: %s", err)
			}
		})
	}
}
