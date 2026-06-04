package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
)

func (c *service) Get(ctx context.Context, uuid uuid.UUID) (model.Order, error) {
	order, err := c.orderRepo.Get(ctx, uuid)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}
	return order, nil
}
