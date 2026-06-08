package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, rawUUID string) (model.Part, error) {
	_, err := uuid.Parse(rawUUID)
	if err != nil {
		return model.Part{}, errs.ErrInvalidUUID
	}

	return s.partRepo.Get(ctx, rawUUID)
}
