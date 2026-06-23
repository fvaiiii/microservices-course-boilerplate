package order

import (
	"context"
	"fmt"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := r.createOrder(txCtx, order); err != nil {
			return err
		}
		return r.createOrderItems(txCtx, order)
	})
}

func (r *repository) createOrder(ctx context.Context, order model.Order) error {
	rec := converter.OrderToRecord(order)

	const query = `
		INSERT INTO orders (uuid, status, created_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query,
		rec.UUID,
		rec.Status,
		rec.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}

func (r *repository) createOrderItems(ctx context.Context, order model.Order) error {
	if len(order.Items) == 0 {
		return nil
	}

	items := converter.OrderItemsToRecord(order.UUID, order.Items)
	const query = `
		INSERT INTO order_items (order_uuid, part_uuid, part_type, price)
		VALUES ($1, $2, $3, $4)
	`

	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	for _, item := range items {
		_, err := db.Exec(ctx, query,
			item.OrderUUID,
			item.PartUUID,
			item.PartType,
			item.Price,
		)
		if err != nil {
			return fmt.Errorf("создать строки заказа: %w", err)
		}
	}

	return nil
}
