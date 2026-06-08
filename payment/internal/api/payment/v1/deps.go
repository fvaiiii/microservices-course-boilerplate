package v1

import (
	"context"

	"github.com/fvaiiii/microservices-course-boilerplate/payment/internal/service/input"
)

type PaymentService interface {
	Pay(ctx context.Context, in input.PayOrderInput) (string, error)
}
