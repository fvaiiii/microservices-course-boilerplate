package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
)

func (c *service) Pay(ctx context.Context, id uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	order, err := c.orderRepo.Get(ctx, id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("получить заказ: %w", err)
	}
	if order.Status == model.OrderStatusPaid {
		return uuid.Nil, errs.ErrOrderAlreadyPaid
	}
	if order.Status == model.OrderStatusCancelled {
		return uuid.Nil, errs.ErrOrderCancelled
	}

	transactionUUID, err := c.paymentRepo.PayOrder(ctx, id, method)
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}
	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &method

	if err := c.orderRepo.Update(ctx, order); err != nil {
		return uuid.Nil, fmt.Errorf("обновить заказ: %w", err)
	}
	return transactionUUID, nil
}
