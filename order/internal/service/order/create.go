package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/input"
)

func (s *service) Create(ctx context.Context, in input.CreateOrderInput) (model.Order, error) {
	parts, err := s.inventoryRepo.ListParts(ctx, in.PartUUIDs())
	if err != nil {
		return model.Order{}, fmt.Errorf("получить детали: %w", err)
	}

	items := make([]model.OrderItem, 0, len(parts))
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, fmt.Errorf("деталь %s: %w", part.Name, errs.ErrOutOfStock)
		}
		items = append(items, model.OrderItem{
			PartUUID: part.UUID,
			PartType: part.PartType,
			Price:    part.Price,
		})
	}

	order := model.Order{
		UUID:      uuid.New(),
		Items:     items,
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return model.Order{}, fmt.Errorf("создать заказ: %w", err)
	}

	return order, nil
}
