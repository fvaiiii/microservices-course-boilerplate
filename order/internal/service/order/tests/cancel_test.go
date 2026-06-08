package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	ordersvc "github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/order"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/order/mocks"
)

func TestCancel(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = uuid.New()

		pendingOrder = model.Order{
			UUID:   orderUUID,
			Status: model.OrderStatusPendingPayment,
		}
		paidOrder = model.Order{
			UUID:   orderUUID,
			Status: model.OrderStatusPaid,
		}
		cancelledOrder = model.Order{
			UUID:   orderUUID,
			Status: model.OrderStatusCancelled,
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
			name: "успешная отмена",
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(pendingOrder, nil)

				repo.EXPECT().
					Update(ctx, mock.MatchedBy(func(o model.Order) bool {
						return o.UUID == orderUUID && o.Status == model.OrderStatusCancelled
					})).
					Return(nil)
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
			name: "заказ уже оплачен",
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(paidOrder, nil)
			},
			expected: expected{err: errs.ErrOrderAlreadyPaid},
		},
		{
			name: "заказ уже отменён",
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(cancelledOrder, nil)
			},
			expected: expected{err: errs.ErrOrderCancelled},
		},
		{
			name: "ошибка при обновлении заказа",
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(pendingOrder, nil)

				repo.EXPECT().
					Update(ctx, mock.Anything).
					Return(errors.New("ошибка БД"))
			},
			expected: expected{onlyError: true, containsMsg: "обновить заказ"},
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
			err := svc.Cancel(ctx, orderUUID)

			switch {
			case tc.expected.onlyError:
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expected.containsMsg)

			case tc.expected.err != nil:
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)

			default:
				require.NoError(t, err)
			}
		})
	}
}
