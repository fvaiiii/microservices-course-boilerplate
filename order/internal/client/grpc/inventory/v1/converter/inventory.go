package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
)

func ToParts(parts []*inventoryv1.Part) ([]model.Part, error) {
	result := make([]model.Part, 0, len(parts))

	for _, part := range parts {
		id, err := uuid.Parse(part.Uuid)
		if err != nil {
			return nil, fmt.Errorf("parse uuid: %w", err)
		}

		result = append(result, model.Part{
			UUID:          id,
			Name:          part.Name,
			PartType:      ToPartType(part.PartType),
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
		})
	}

	return result, nil
}

func ToPartType(t inventoryv1.PartType) model.PartType {
	switch t {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull

	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine

	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield

	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon

	default:
		return ""
	}
}
