package v1

import (
	"context"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/input"
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	input := input.PayOrderInput{
		OrderUUID: params.OrderUUID,
		Method:    model.PaymentMethod(req.PaymentMethod),
	}
	transactionUUID, err := a.orderService.Pay(ctx, input.OrderUUID, input.Method)
	if err != nil {
		return handlePayOrderError(err), nil
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
