package part

import (
	"context"

	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	List(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
}
