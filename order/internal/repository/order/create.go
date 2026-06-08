package order

import (
	"context"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[order.UUID] = converter.OrderToRecord(order)

	return nil
}
