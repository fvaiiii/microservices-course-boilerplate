package converter

import (
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/repository/record"
)

func PartToModel(part record.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      model.PartType(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}

func PartToRecord(part model.Part) record.Part {
	return record.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      string(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}
