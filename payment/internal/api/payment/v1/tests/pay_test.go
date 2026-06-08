package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	paymentapi "github.com/fvaiiii/microservices-course-boilerplate/payment/internal/api/payment/v1"
	"github.com/fvaiiii/microservices-course-boilerplate/payment/internal/api/payment/v1/mocks"
	errs "github.com/fvaiiii/microservices-course-boilerplate/payment/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/payment/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/payment/internal/service/input"
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		orderUUID       = uuid.New()
		transactionUUID = uuid.New()
	)

	tests := []struct {
		name      string
		req       *paymentv1.PayOrderRequest
		setupMock func(scv *mocks.PaymentService)
		wantErr   error
	}{
		{
			name: "успешная оплата картой",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     orderUUID.String(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					Pay(ctx, input.PayOrderInput{
						OrderUUID:     orderUUID.String(),
						PaymentMethod: model.PaymentMethodCard,
					}).
					Return(transactionUUID.String(), nil)
			},
		},
		{
			name: "успешная оплата через СБП",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     orderUUID.String(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
			},
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					Pay(ctx, input.PayOrderInput{
						OrderUUID:     orderUUID.String(),
						PaymentMethod: model.PaymentMethodSBP,
					}).
					Return(transactionUUID.String(), nil)
			},
		},
		{
			name: "неверный метод оплаты",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     orderUUID.String(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
			},
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					Pay(ctx, input.PayOrderInput{
						OrderUUID:     orderUUID.String(),
						PaymentMethod: model.PaymentMethodUnspecified,
					}).
					Return("", errs.ErrInvalidPaymentMethod)
			},
			wantErr: errs.ErrInvalidPaymentMethod,
		},
		{
			name: "неверный UUID заказа",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     "",
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					Pay(ctx, input.PayOrderInput{
						OrderUUID:     "",
						PaymentMethod: model.PaymentMethodCard,
					}).
					Return("", errs.ErrInvalidOrderUUID)
			},
			wantErr: errs.ErrInvalidOrderUUID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewPaymentService(t)
			tc.setupMock(svc)

			api := paymentapi.New(svc)
			res, err := api.PayOrder(ctx, tc.req)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
			assert.Equal(t, transactionUUID.String(), res.GetTransactionUuid())
		})
	}
}
