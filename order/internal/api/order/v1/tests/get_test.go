package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderapi "github.com/fvaiiii/microservices-course-boilerplate/order/internal/api/order/v1"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/api/order/v1/mocks"
	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

func TestGetOrder(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		orderUUID  = uuid.New()
		hullUUID   = uuid.New()
		engineUUID = uuid.New()
		createdAt  = time.Now()
	)

	tests := []struct {
		name       string
		setupMock  func(svc *mocks.OrderService)
		wantStatus int
	}{
		{
			name: "успешное получение заказа",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						UUID: orderUUID,
						Items: []model.OrderItem{
							{PartUUID: hullUUID, PartType: model.PartTypeHull, Price: 500000},
							{PartUUID: engineUUID, PartType: model.PartTypeEngine, Price: 300000},
						},
						Status:    model.OrderStatusPendingPayment,
						CreatedAt: createdAt,
					}, nil)
			},
			wantStatus: 200,
		},
		{
			name: "заказ не найден",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			wantStatus: 404,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewOrderService(t)
			tc.setupMock(svc)

			api := orderapi.New(svc)
			res, err := api.GetOrder(ctx, orderv1.GetOrderParams{OrderUUID: orderUUID})

			require.NoError(t, err)

			switch tc.wantStatus {
			case 200:
				dto, ok := res.(*orderv1.OrderDto)
				require.True(t, ok)
				assert.Equal(t, orderUUID, dto.OrderUUID)
				assert.Equal(t, hullUUID, dto.HullUUID)
				assert.Equal(t, engineUUID, dto.EngineUUID)
				assert.Equal(t, int64(800000), dto.TotalPrice)
				assert.Equal(t, orderv1.OrderStatusPENDINGPAYMENT, dto.Status)

			case 404:
				resp, ok := res.(*orderv1.GetOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, 404, resp.Code)
			}
		})
	}
}
