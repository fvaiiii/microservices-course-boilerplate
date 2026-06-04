package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, uuid uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.data[uuid]
	if !exists {
		return model.Order{}, errs.ErrOrderNotFound
	}
	return converter.RecordToModel(order), nil
}
