package order

import (
	"context"
	"fmt"
	"time"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		return r.updateOrder(txCtx, order)
	})
}

func (r *repository) updateOrder(ctx context.Context, order model.Order) error {
	rec := converter.OrderToRecord(order)
	const query = `
		UPDATE orders
		SET status = $1,
		    transaction_uuid = $2,
		    payment_method = $3,
		    updated_at = $4
		WHERE uuid = $5
	`
	tag, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query,
		rec.Status,
		rec.TransactionUUID,
		rec.PaymentMethod,
		time.Now(),
		rec.UUID,
	)
	if err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}
	return nil
}
