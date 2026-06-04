package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/input"
)

type OrderService interface {
	Create(ctx context.Context, in input.CreateOrderInput) (model.Order, error)
	Get(ctx context.Context, uuid uuid.UUID) (model.Order, error)
	Pay(ctx context.Context, uuid uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
	Cancel(ctx context.Context, uuid uuid.UUID) error
}
