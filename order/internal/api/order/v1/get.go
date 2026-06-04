package v1

import (
	"context"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/api/converter"
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		return handleGetOrderError(err), nil
	}

	return converter.ToOrderDTO(order), nil
}
