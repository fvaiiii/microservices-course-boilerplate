package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
)

func (c *service) Cancel(ctx context.Context, uuid uuid.UUID) error {
	order, err := c.orderRepo.Get(ctx, uuid)
	if err != nil {
		return fmt.Errorf("получить заказ: %w", err)
	}
	if order.Status == model.OrderStatusPaid {
		return errs.ErrOrderAlreadyPaid
	}

	if order.Status == model.OrderStatusCancelled {
		return errs.ErrOrderCancelled
	}

	order.Status = model.OrderStatusCancelled
	if err := c.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}

	return nil
}
