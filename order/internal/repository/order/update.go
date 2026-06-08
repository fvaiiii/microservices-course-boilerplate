package order

import (
	"context"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[order.UUID]; !exists {
		return errs.ErrOrderNotFound
	}

	r.data[order.UUID] = converter.OrderToRecord(order)
	return nil
}
