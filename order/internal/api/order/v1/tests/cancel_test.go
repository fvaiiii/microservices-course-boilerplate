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
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

func TestCancelOrder(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
	)

	tests := []struct {
		name       string
		setupMock  func(svc *mocks.OrderService)
		wantStatus int
	}{
		{
			name: "успешная отмена",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID).
					Return(nil)
			},
			wantStatus: 200,
		},
		{
			name: "заказ не найден",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID).
					Return(errs.ErrOrderNotFound)
			},
			wantStatus: 404,
		},
		{
			name: "заказ уже оплачен",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID).
					Return(errs.ErrOrderAlreadyPaid)
			},
			wantStatus: 409,
		},
		{
			name: "заказ уже отменён",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID).
					Return(errs.ErrOrderCancelled)
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
			res, err := api.CancelOrder(ctx, orderv1.CancelOrderParams{OrderUUID: orderUUID})

			require.NoError(t, err)

			switch tc.wantStatus {
			case 200:
				_, ok := res.(*orderv1.CancelOrderResponse)
				require.True(t, ok)

			case 404:
				resp, ok := res.(*orderv1.CancelOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, 404, resp.Code)

			case 409:
				resp, ok := res.(*orderv1.CancelOrderConflict)
				require.True(t, ok)
				assert.Equal(t, 409, resp.Code)
			}
		})
	}
}
