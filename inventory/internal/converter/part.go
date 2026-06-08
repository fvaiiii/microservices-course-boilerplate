package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
)

func PartToProto(part model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      PartTypeToProto(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     timestamppb.New(part.CreatedAt),
	}
}

func PartTypeToModel(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return "HULL"
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return "ENGINE"
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return "SHIELD"
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return "WEAPON"
	default:
		return "UNSPECIFIED"
	}
}

func PartTypeToProto(partType model.PartType) inventoryv1.PartType {
	switch partType {
	case model.PartTypeHull:
		return inventoryv1.PartType_PART_TYPE_HULL
	case model.PartTypeEngine:
		return inventoryv1.PartType_PART_TYPE_ENGINE
	case model.PartTypeShield:
		return inventoryv1.PartType_PART_TYPE_SHIELD
	case model.PartTypeWeapon:
		return inventoryv1.PartType_PART_TYPE_WEAPON
	default:
		return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
	}
}
