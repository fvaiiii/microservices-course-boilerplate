package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/converter"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, id uuid.UUID) (model.Order, error) {
	orderRec, err := r.getOrder(ctx, id)
	if err != nil {
		return model.Order{}, err
	}
	items, err := r.getOrderItems(ctx, id)
	if err != nil {
		return model.Order{}, err
	}

	return converter.RecordToModel(orderRec, items), nil
}

func (r *repository) getOrder(ctx context.Context, id uuid.UUID) (record.Order, error) {
	const query = `
		SELECT uuid, status, transaction_uuid, payment_method, created_at, updated_at
		FROM orders
		WHERE uuid = $1
	`
	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, id)
	if err != nil {
		return record.Order{}, fmt.Errorf("получить заказ: %w", err)
	}
	defer rows.Close()

	orderRec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[record.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return record.Order{}, errs.ErrOrderNotFound
		}
		return record.Order{}, fmt.Errorf("получить заказ: %w", err)
	}
	return orderRec, nil
}

func (r *repository) getOrderItems(ctx context.Context, orderUUID uuid.UUID) ([]record.OrderItem, error) {
	const query = `
		SELECT order_uuid, part_uuid, part_type, price
		FROM order_items
		WHERE order_uuid = $1
	`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, orderUUID)
	if err != nil {
		return nil, fmt.Errorf("получить строки заказа: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.OrderItem])
	if err != nil {
		return nil, fmt.Errorf("получить строки заказа: %w", err)
	}

	return items, nil
}
