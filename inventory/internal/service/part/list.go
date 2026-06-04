package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	for _, rawUUID := range filter.UUIDs {
		if _, err := uuid.Parse(rawUUID); err != nil {
			return nil, errs.ErrInvalidUUID
		}
	}

	return s.partRepo.List(ctx, filter)
}
