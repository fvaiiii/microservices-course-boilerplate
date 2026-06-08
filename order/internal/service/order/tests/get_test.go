package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	ordersvc "github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/order"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/order/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	var (
		orderUUID     = uuid.New()
		ctx           = context.Background()
		existingOrder = model.Order{
			UUID:   orderUUID,
			Status: model.OrderStatusPendingPayment,
		}
	)

	type expected struct {
		err         error
		onlyError   bool
		containsMsg string
	}

	tests := []struct {
		name      string
		setupMock func(repo *mocks.OrderRepository)
		expected  expected
	}{
		{
			name: "успешное получение заказа",
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(existingOrder, nil)
			},
			expected: expected{err: nil},
		},
		{
			name: "заказ не найден",
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{err: errs.ErrOrderNotFound},
		},
		{
			name: "ошибка репозитория",
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, errors.New("ошибка БД"))
			},
			expected: expected{onlyError: true, containsMsg: "получить заказ"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			orderRepo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)
			tc.setupMock(orderRepo)
			svc := ordersvc.New(orderRepo, inventoryClient, paymentClient)
			order, err := svc.Get(ctx, orderUUID)
			switch {
			case tc.expected.onlyError:
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expected.containsMsg)
				assert.Empty(t, order.UUID)
			case tc.expected.err != nil:
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Empty(t, order.UUID)
			default:
				require.NoError(t, err)
				assert.Equal(t, orderUUID, order.UUID)
				assert.Equal(t, model.OrderStatusPendingPayment, order.Status)
			}
		})
	}
}
