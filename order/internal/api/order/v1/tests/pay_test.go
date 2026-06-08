package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderapi "github.com/fvaiiii/microservices-course-boilerplate/order/internal/api/order/v1"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/api/order/v1/mocks"
	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		orderUUID       = uuid.New()
		transactionUUID = uuid.New()
	)

	tests := []struct {
		name       string
		setupMock  func(svc *mocks.OrderService)
		wantStatus int
	}{
		{
			name: "успешная оплата",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, model.PaymentMethodCard).
					Return(transactionUUID, nil)
			},
			wantStatus: 200,
		},
		{
			name: "заказ не найден",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, model.PaymentMethodCard).
					Return(uuid.Nil, errs.ErrOrderNotFound)
			},
			wantStatus: 404,
		},
		{
			name: "заказ уже оплачен",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, model.PaymentMethodCard).
					Return(uuid.Nil, errs.ErrOrderAlreadyPaid)
			},
			wantStatus: 409,
		},
		{
			name: "заказ отменён",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, model.PaymentMethodCard).
					Return(uuid.Nil, errs.ErrOrderCancelled)
			},
			wantStatus: 409,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewOrderService(t)
			tc.setupMock(svc)

			api := orderapi.New(svc)
			res, err := api.PayOrder(
				ctx,
				&orderv1.PayOrderRequest{PaymentMethod: orderv1.PaymentMethodCARD},
				orderv1.PayOrderParams{OrderUUID: orderUUID},
			)

			require.NoError(t, err)

			switch tc.wantStatus {
			case 200:
				resp, ok := res.(*orderv1.PayOrderResponse)
				require.True(t, ok)
				assert.Equal(t, transactionUUID, resp.TransactionUUID)

			case 404:
				resp, ok := res.(*orderv1.PayOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, 404, resp.Code)

			case 409:
				resp, ok := res.(*orderv1.PayOrderConflict)
				require.True(t, ok)
				assert.Equal(t, 409, resp.Code)
			}
		})
	}
}
