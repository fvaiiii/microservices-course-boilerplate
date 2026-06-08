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

func TestPay(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		orderUUID       = uuid.New()
		transactionUUID = uuid.New()
		paymentMethod   = model.PaymentMethodCard

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
		transactionUUID uuid.UUID
		err             error
		onlyError       bool
		containsMsg     string
	}

	tests := []struct {
		name      string
		setupMock func(repo *mocks.OrderRepository, payment *mocks.PaymentClient)
		expected  expected
	}{
		{
			name: "успешная оплата",
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(pendingOrder, nil)

				payment.EXPECT().
					PayOrder(ctx, orderUUID, paymentMethod).
					Return(transactionUUID, nil)

				repo.EXPECT().
					Update(ctx, mock.MatchedBy(func(o model.Order) bool {
						return o.UUID == orderUUID &&
							o.Status == model.OrderStatusPaid &&
							o.TransactionUUID != nil &&
							*o.TransactionUUID == transactionUUID &&
							o.PaymentMethod != nil &&
							*o.PaymentMethod == paymentMethod
					})).
					Return(nil)
			},
			expected: expected{transactionUUID: transactionUUID, err: nil},
		},
		{
			name: "заказ не найден",
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{err: errs.ErrOrderNotFound},
		},
		{
			name: "заказ уже оплачен",
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(paidOrder, nil)
			},
			expected: expected{err: errs.ErrOrderAlreadyPaid},
		},
		{
			name: "заказ отменён",
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(cancelledOrder, nil)
			},
			expected: expected{err: errs.ErrOrderCancelled},
		},
		{
			name: "ошибка PaymentService",
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(pendingOrder, nil)

				payment.EXPECT().
					PayOrder(ctx, orderUUID, paymentMethod).
					Return(uuid.Nil, errors.New("ошибка оплаты"))
			},
			expected: expected{onlyError: true, containsMsg: "оплатить заказ"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)

			tc.setupMock(orderRepo, paymentClient)

			svc := ordersvc.New(orderRepo, inventoryClient, paymentClient)
			txUUID, err := svc.Pay(ctx, orderUUID, paymentMethod)

			switch {
			case tc.expected.onlyError:
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expected.containsMsg)
				assert.Equal(t, uuid.Nil, txUUID)

			case tc.expected.err != nil:
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, uuid.Nil, txUUID)

			default:
				require.NoError(t, err)
				assert.Equal(t, tc.expected.transactionUUID, txUUID)
			}
		})
	}
}
