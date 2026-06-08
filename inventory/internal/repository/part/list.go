package part

import (
	"context"
	"sort"

	"github.com/google/uuid"

	errs "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
	repoConverter "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/repository/converter"
)

func (r *inventoryRepo) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(filter.UUIDs) > 0 {
		partsArr := make([]model.Part, 0, len(filter.UUIDs))

		for _, rawUUID := range filter.UUIDs {
			id, err := uuid.Parse(rawUUID)
			if err != nil {
				return nil, errs.ErrInvalidUUID
			}

			part, ok := r.parts[id]
			if !ok {
				return nil, errs.ErrPartNotFound
			}

			partsArr = append(partsArr, repoConverter.PartToModel(part))

		}
		return partsArr, nil
	}

	partsArr := make([]model.Part, 0, len(r.parts))
	for _, part := range r.parts {
		if filter.PartType == model.PartTypeUnspecified || part.PartType == string(filter.PartType) {
			partsArr = append(partsArr, repoConverter.PartToModel(part))
		}
	}

	sort.Slice(partsArr, func(i, j int) bool {
		return partsArr[i].Name < partsArr[j].Name
	})
	return partsArr, nil
}
