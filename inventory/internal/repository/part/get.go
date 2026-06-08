package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
	repoConverter "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/repository/converter"
)

func (r *inventoryRepo) Get(ctx context.Context, rawUUID string) (model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, err := uuid.Parse(rawUUID)
	if err != nil {
		return model.Part{}, errs.ErrInvalidUUID
	}
	part, ok := r.parts[id]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}

	return repoConverter.PartToModel(part), nil
}
