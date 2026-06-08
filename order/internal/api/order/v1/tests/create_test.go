package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	orderapi "github.com/fvaiiii/microservices-course-boilerplate/order/internal/api/order/v1"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/api/order/v1/mocks"
	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/input"
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		hullUUID   = uuid.New()
		engineUUID = uuid.New()
		orderUUID  = uuid.New()
	)

	type expected struct {
		statusCode int
		orderUUID  uuid.UUID
		totalPrice int64
	}

	tests := []struct {
		name      string
		req       *orderv1.CreateOrderRequest
		setupMock func(svc *mocks.OrderService)
		expected  expected
	}{
		{
			name: "успешное создание заказа",
			req: &orderv1.CreateOrderRequest{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, input.CreateOrderInput{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
					}).
					Return(model.Order{
						UUID: orderUUID,
						Items: []model.OrderItem{
							{PartUUID: hullUUID, Price: 500000},
							{PartUUID: engineUUID, Price: 300000},
						},
						Status: model.OrderStatusPendingPayment,
					}, nil)
			},
			expected: expected{
				statusCode: 201,
				orderUUID:  orderUUID,
				totalPrice: 800000,
			},
		},
		{
			name: "деталь не найдена",
			req: &orderv1.CreateOrderRequest{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, mock.Anything).
					Return(model.Order{}, errs.ErrPartNotFound)
			},
			expected: expected{statusCode: 404},
		},
		{
			name: "деталь отсутствует на складе",
			req: &orderv1.CreateOrderRequest{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, mock.Anything).
					Return(model.Order{}, errs.ErrOutOfStock)
			},
			expected: expected{statusCode: 409},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewOrderService(t)
			tc.setupMock(svc)

			api := orderapi.New(svc)
			res, err := api.CreateOrder(ctx, tc.req)

			require.NoError(t, err)

			switch tc.expected.statusCode {
			case 201:
				resp, ok := res.(*orderv1.CreateOrderResponse)
				require.True(t, ok)
				assert.Equal(t, tc.expected.orderUUID, resp.OrderUUID)
				assert.Equal(t, tc.expected.totalPrice, resp.TotalPrice)

			case 404:
				resp, ok := res.(*orderv1.CreateOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, 404, resp.Code)

			case 409:
				resp, ok := res.(*orderv1.CreateOrderConflict)
				require.True(t, ok)
				assert.Equal(t, 409, resp.Code)
			}
		})
	}
}
