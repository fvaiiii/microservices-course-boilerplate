package v1

import (
	"context"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/api/converter"
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	in := converter.ToCreateOrderInput(req)
	order, err := a.orderService.Create(ctx, in)
	if err != nil {
		return handleCreateOrderError(err), nil
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice(),
	}, nil
}
